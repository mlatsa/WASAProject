from fastapi import FastAPI, Header, HTTPException, Depends, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse, PlainTextResponse
from pydantic import BaseModel, Field, HttpUrl
from typing import Dict, List, Optional
import uuid
import time

app = FastAPI(title="WASAText API")


app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.options("/{full_path:path}")
async def cors_preflight(full_path: str):
    resp = PlainTextResponse("")
    resp.headers["Access-Control-Allow-Origin"] = "*"
    resp.headers["Access-Control-Allow-Methods"] = "GET,POST,PUT,DELETE,OPTIONS,PATCH"
    resp.headers["Access-Control-Allow-Headers"] = "Authorization,Content-Type"
    resp.headers["Access-Control-Max-Age"] = "1"
    return resp

# --- In-memory stores ---
users_by_name: Dict[str, str] = {}            # name -> identifier
usernames_by_id: Dict[str, str] = {}          # id -> name
user_photos: Dict[str, str] = {}              # id -> mediaUrl

class Message(BaseModel):
    messageId: str
    conversationId: str
    sender: str
    content: str
    type: str = "text"
    status: str = "delivered"
    timestamp: float = Field(default_factory=lambda: time.time())

class Conversation(BaseModel):
    id: str
    participants: List[str]
    lastMessage: Optional[str] = None
    timestamp: float = Field(default_factory=lambda: time.time())

conversations: Dict[str, Conversation] = {}
messages_by_convo: Dict[str, List[Message]] = {}
reactions_by_message: Dict[str, Dict[str, str]] = {}

def require_bearer(authorization: Optional[str] = Header(None)) -> str:
    if not authorization or not authorization.lower().startswith("bearer "):
        raise HTTPException(status_code=401, detail="Unauthorized")
    token = authorization.split(" ", 1)[1].strip()
    if token not in usernames_by_id:
        raise HTTPException(status_code=401, detail="Unauthorized")
    return token  # return identifier

# --- Schemas (requests/responses) ---
class LoginBody(BaseModel):
    name: str = Field(min_length=3, max_length=16)

class SetNameBody(BaseModel):
    name: str = Field(min_length=3, max_length=16, pattern=r"^[a-zA-Z0-9_.-]{3,16}$")

class SendMessageInput(BaseModel):
    content: str
    type: str = "text"

class ForwardBody(BaseModel):
    conversationId: str

class ReactionBody(BaseModel):
    reaction: str

class PhotoBody(BaseModel):
    mediaUrl: HttpUrl

class GroupAddBody(BaseModel):
    id: str

class GroupNameBody(BaseModel):
    name: str = Field(min_length=3, max_length=32)

# --- Utilities ---
def ensure_conversation_exists(conversation_id: str, participants: Optional[List[str]] = None):
    if conversation_id not in conversations:
        if not participants:
            raise HTTPException(status_code=404, detail="Conversation not found")
        conversations[conversation_id] = Conversation(id=conversation_id, participants=participants)
        messages_by_convo[conversation_id] = []

def new_id(prefix: str) -> str:
    return f"{prefix}_{uuid.uuid4().hex[:12]}"

# --- Endpoints ---

# Health 
@app.get("/health")
def health():
    return {"status": "ok"}

# session 
@app.post("/session", status_code=201)
def do_login(body: LoginBody):
    name = body.name
    if name in users_by_name:
        identifier = users_by_name[name]
    else:
        identifier = new_id("user")
        users_by_name[name] = identifier
        usernames_by_id[identifier] = name
        
    return {"identifier": identifier}

# /user/username (setMyUserName)
@app.post("/user/username")
def set_my_username(body: SetNameBody, user_id: str = Depends(require_bearer)):
    new_name = body.name
    # is it taken by someone else?
    if new_name in users_by_name and users_by_name[new_name] != user_id:
        raise HTTPException(status_code=400, detail="Invalid username.")
    # remove old name mapping
    old_name = usernames_by_id.get(user_id)
    if old_name and old_name in users_by_name:
        del users_by_name[old_name]
    # set new
    users_by_name[new_name] = user_id
    usernames_by_id[user_id] = new_name
    return {"message": "Username updated."}

# /conversations (getMyConversations)
@app.get("/conversations")
def get_my_conversations(user_id: str = Depends(require_bearer)):
    # Return all conversations where user participates
    mine = [c for c in conversations.values() if user_id in c.participants]
    return {"conversations": [c.model_dump() for c in mine]}


@app.get("/conversations/{conversationId}")
def get_conversation(conversationId: str, user_id: str = Depends(require_bearer)):
    if conversationId not in conversations:
        raise HTTPException(status_code=404, detail="Conversation not found")
    conv = conversations[conversationId]
    if user_id not in conv.participants:
        raise HTTPException(status_code=404, detail="Conversation not found")
    msgs = messages_by_convo.get(conversationId, [])
    return {
        "conversation": {
            **conv.model_dump(),
            "messages": [m.model_dump() for m in msgs],
        }
    }


@app.post("/conversations/{conversationId}/messages", status_code=201)
def send_message(conversationId: str, body: SendMessageInput, user_id: str = Depends(require_bearer)):
    # check conversation exists; if not, create a 1:1 with the sender only (or enforce 404)
    ensure_conversation_exists(conversationId, participants=[user_id])
    msg = Message(
        messageId=new_id("msg"),
        conversationId=conversationId,
        sender=user_id,
        content=body.content,
        type=body.type or "text",
        status="delivered",
    )
    messages_by_convo[conversationId].append(msg)
    conv = conversations[conversationId]
    conv.lastMessage = body.content
    conv.timestamp = time.time()
    return msg


@app.post("/messages/{messageId}/forward")
def forward_message(messageId: str, body: ForwardBody, user_id: str = Depends(require_bearer)):
    
    src = None
    for msgs in messages_by_convo.values():
        for m in msgs:
            if m.messageId == messageId:
                src = m
                break
        if src:
            break
    if not src:
        raise HTTPException(status_code=404, detail="Message not found")
    # check destination conversation exists
    ensure_conversation_exists(body.conversationId, participants=[user_id])
    fwd = Message(
        messageId=new_id("msg"),
        conversationId=body.conversationId,
        sender=user_id,
        content=src.content,
        type=src.type,
        status="delivered",
    )
    messages_by_convo[body.conversationId].append(fwd)
    conversations[body.conversationId].lastMessage = fwd.content
    conversations[body.conversationId].timestamp = time.time()
    return fwd

# /messages/{messageId}/reactions
@app.post("/messages/{messageId}/reactions")
def add_reaction(messageId: str, body: ReactionBody, user_id: str = Depends(require_bearer)):
    
    found = False
    for msgs in messages_by_convo.values():
        for m in msgs:
            if m.messageId == messageId:
                found = True
                break
        if found:
            break
    if not found:
        raise HTTPException(status_code=404, detail="Message not found")
    react_id = new_id("react")
    reactions_by_message.setdefault(messageId, {})[react_id] = body.reaction
    return {"reactionId": react_id}


@app.delete("/messages/{messageId}/reactions/{reactionId}", status_code=204)
def remove_reaction(messageId: str, reactionId: str, user_id: str = Depends(require_bearer)):
    store = reactions_by_message.get(messageId)
    if not store or reactionId not in store:
        raise HTTPException(status_code=404, detail="Reaction not found")
    del store[reactionId]
    return JSONResponse(status_code=204, content=None)


@app.delete("/messages/{messageId}", status_code=204)
def delete_message(messageId: str, user_id: str = Depends(require_bearer)):
    for conv_id, msgs in messages_by_convo.items():
        for i, m in enumerate(msgs):
            if m.messageId == messageId:
                del msgs[i]
                return JSONResponse(status_code=204, content=None)
    raise HTTPException(status_code=404, detail="Message not found")


@app.post("/groups/{conversationId}/members")
def add_to_group(conversationId: str, body: GroupAddBody, user_id: str = Depends(require_bearer)):
    ensure_conversation_exists(conversationId, participants=[user_id])
    conv = conversations[conversationId]
    if body.id not in conv.participants:
        conv.participants.append(body.id)
    return {"message": "OK"}

@app.post("/groups/{conversationId}/leave")
def leave_group(conversationId: str, user_id: str = Depends(require_bearer)):
    if conversationId not in conversations:
        raise HTTPException(status_code=404, detail="Conversation not found")
    conv = conversations[conversationId]
    if user_id in conv.participants:
        conv.participants.remove(user_id)
    return {"message": "OK"}

@app.post("/groups/{conversationId}/name")
def set_group_name(conversationId: str, body: GroupNameBody, user_id: str = Depends(require_bearer)):
    ensure_conversation_exists(conversationId, participants=[user_id])
    # we don’t keep a name field in Conversation schema above; add it ad-hoc:
    conv = conversations[conversationId]
    setattr(conv, "name", body.name)
    return {"message": "OK"}

@app.post("/user/photo")
def set_my_photo(body: PhotoBody, user_id: str = Depends(require_bearer)):
    user_photos[user_id] = str(body.mediaUrl)
    return {"message": "OK"}

@app.post("/groups/{conversationId}/photo")
def set_group_photo(conversationId: str, body: PhotoBody, user_id: str = Depends(require_bearer)):
    ensure_conversation_exists(conversationId, participants=[user_id])
    # store on the Conversation object dynamically
    conv = conversations[conversationId]
    setattr(conv, "photo", str(body.mediaUrl))
    return {"message": "OK"}


