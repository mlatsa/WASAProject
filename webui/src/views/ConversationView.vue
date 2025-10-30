<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	emits: ['selectedMessageChanged', 'conversationUpdated', 'reloginNeeded'],
	props: ['session_token', 'user', 'selected_conversation', 'need_update_conversation'],
	data: function() {
		return {
			errormsg: '',
			messages: [],
			// selected_message_id: 0
		}
	},
	methods: {
		async doLoadConversation() {
			try {
				if (this.selected_conversation == null)
				{
					this.messages = [];
				} else {
					const response = await this.$axios({
						method: 'get',
						url: '/conversations/' + this.selected_conversation.conversation_id,
						headers: {
							'Authorization' : 'Bearer ' + this.session_token
						}
					});

					if (response.status != 200) {
						this.errormsg = "Unexpected response " + response.status;
						return;
					}

					this.messages = response.data;
				}
				this.errormsg = "";
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
		},
		async onMessageClick(msg) {
			// this.selected_message_id = msg.msg_id;
			this.$emit('selectedMessageChanged', msg);
		},
		async onErrorDismissed() {
			this.errormsg = '';
		}
	},
	watch: {
    	session_token(newValue, oldValue) {
			if (newValue == 0) {
				this.errormsg = '';
				this.messages = [];
				// this.selected_message_id = 0;
			}
		},
    	selected_conversation(newValue, oldValue) {
      		this.doLoadConversation();
		},
    	need_update_conversation(newValue, oldValue) {
      		this.doLoadConversation();
			this.$emit('conversationUpdated');
		}
	},
	beforeMount: function () {
        window.addEventListener('beforeunload', (e) => {
			try {
				// non salva i messaggi perché sarebbe troppo pesante, occorre riaprire la conversazione
            	// sessionStorage.setItem('messages',  JSON.stringify(this.messages));
            	// sessionStorage.setItem('selected_message_id', JSON.stringify(this.selected_message_id));
			} catch {
				this.$emit('reloginNeeded');
			}
        });

        try {
            // if (sessionStorage.getItem('messages') != null) this.messages = JSON.parse(sessionStorage.getItem('messages'));
            // if (sessionStorage.getItem('selected_message_id') != null) this.selected_message_id = JSON.parse(sessionStorage.getItem('selected_message_id'));
        } catch {
			this.$emit('reloginNeeded');
        }
    },
	mounted: function () {
		this.doLoadConversation();
	}
}
</script>

<template>
	<div v-if="session_token != 0">
		<div ref="convview" v-if="session_token != 0 && user != null && selected_conversation != null" class="conversation-container">
			<div :class="{'message-box friend-message':msg.from_user_id!=user.id,'message-box my-message':msg.from_user_id==user.id}"
				v-for="(msg, index) in messages" :key="msg.id" :tabindex="index" @click="onMessageClick(msg)">
				<p v-if=!msg.is_photo>
					<span style="font-style: italic;">{{ msg.from_user_name }}</span>
					<span style="font-style: italic;">{{ msg.forwarded_from_msg_id > 0 ? '(forwarded)' : '' }}</span>
					<span>{{ msg.msg }}</span>
					<span style="position:relative; font-size: smaller;">
						<span v-if="msg.from_user_id==user.id">{{ msg.sent_timestamp }}&nbsp;{{ msg.seen > 0 ? "&#x2713;&#x2713;" : "&#x2713;" }}</span>
						<span v-if="msg.from_user_id!=user.id">{{ msg.sent_timestamp }}</span>
					</span>
					<span style="font-size: medium">{{ msg.reaction }}</span>
				</p>
				<p v-if=msg.is_photo>
					<span style="font-style: italic;">{{ msg.from_user_name }}</span>
					<span style="font-style: italic;">{{ msg.forwarded_from_msg_id > 0 ? '(forwarded)' : '' }}</span>
					<img class="photo-box-big" v-bind:src=msg.msg>
					<span style="position:relative; font-size: smaller;">
						<span v-if="msg.from_user_id==user.id">{{ msg.sent_timestamp }}&nbsp;{{ msg.seen > 0 ? "&#x2713;&#x2713;" : "&#x2713;" }}</span>
						<span v-if="msg.from_user_id!=user.id">{{ msg.sent_timestamp }}</span>
					</span>
					<span style="font-size: medium">{{ msg.reaction }}</span>
				</p>
			</div>
		</div>
		<ErrorMsg :errormsg="errormsg" @errorDismissed="onErrorDismissed"></ErrorMsg>
	</div>
</template>
