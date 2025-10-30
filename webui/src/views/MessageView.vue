<script setup>
//import { ref } from "vue";
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	emits: ['messageSent', 'messageModified', 'messageForwarded', 'reloginNeeded'],
    props: ['session_token', 'user', 'all_users_and_my_groups', 'selected_message', 'selected_conversation'],
    data: function() {
		return {
			errormsg: '',
			user_or_group_for_send_msg: null,
			user_or_group_for_forward_msg: null,
			message: ''
		}
	},
	methods: {
		async doSendMessage() {
			if (this.selected_conversation == null) {
				return;
			}
            try {
                if (this.message == '') {
                    this.errormsg = "Message not set";
                    return;
                }
                
                const res = await this.$axios({
                    method: 'post',
                    url: '/messages',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "to_user_name_or_group_name": this.selected_conversation.user_or_group_name,
                        "is_group": this.selected_conversation.is_group == 0 ? false : true,
                        "message": this.message
                    }
                });

                if (res.status != 201) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                this.$emit('messageSent', this.selected_conversation.user_or_group_name);
                this.message = '';

                this.errormsg = '';
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        // this code cannot be placed inside doSetMyPhoto because
        // await need to be on the top level of an async method
        async doSendPhotoAsync(img)
        {
            try {
                const res = await this.$axios({
                    method: 'post',
                    url: '/messages',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "to_user_name_or_group_name": this.selected_conversation.user_or_group_name,
                        "is_group": this.selected_conversation.is_group == 0 ? false : true,
                        "message": img
                    }
                });

                if (res.status != 201) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                this.$emit('messageSent', this.for_user_name_or_group_name);

                this.errormsg = '';
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        async doSendPhoto(e) {
			if (this.selected_conversation == null) {
				return;
			}

            const img = e.target.files[0];

            if (img == null)
                return;
            const reader = new FileReader();
            reader.readAsDataURL(img);
            reader.onload = evn =>{
                this.doSendPhotoAsync(evn.target.result);
            }
        },

        async onImageUploaderClick() {
            this.$refs.photoSender.value = ''
        },

		async doForwardMessage() {
            if (this.selected_message == null) {
				return;
			}
			if (this.user_or_group_for_forward_msg == null) {
				this.errormsg = "Receiver user not set";
				return;
			}

			try {
				const response = await this.$axios({
					method: 'post',
					url: '/messages/' + this.selected_message.msg_id,
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					},
					data: {
                        "to_user_name_or_group_name": this.user_or_group_for_forward_msg.user_name_or_group_name,
                        "is_group": this.user_or_group_for_forward_msg.is_group
                    }
				});

                this.$emit('messageForwarded', name);

				this.errormsg = '';
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
		},
        async onDeleteMessage() {
            if (this.selected_message == null) {
				return;
			}
			try {
				const response = await this.$axios({
					method: 'delete',
					url: '/messages/' + this.selected_message.msg_id,
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					}
				});

				this.$emit('messageModified', this.to_user_name);
				this.errormsg = '';
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
        },
		async onDeleteReaction() {
            if (this.selected_message == null) {
				return;
			}
			try {
				const response = await this.$axios({
					method: 'delete',
					url: '/messages/:' + this.selected_message.msg_id + '/comments',
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					}
				});

				this.$emit('messageModified', this.to_user_name);
				this.errormsg = '';
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
		},
		async onReactionClick(el) {
			if (this.selected_message == null) {
				return;
			}

			if (el.target == null|| el.target.value == '') {
				return;
			}
			try {
				const response = await this.$axios({
					method: 'post',
					url: '/messages/:msg_id/comments',
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					},
                    data: {
                        "msg_id": this.selected_message.msg_id,
                        "reaction": el.target.value
                    }
				});

				if (response.status != 201) {
                    this.errormsg = "Unexpected response " + response.status;
                    return;
                }

				this.$emit('messageModified', this.to_user_name);
				this.errormsg = '';
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
		},
		async onErrorDismissed() {
			this.errormsg = '';
		}
	},
	watch: {
		session_token(newValue, oldValue) {
			if (newValue == 0) {
				this.errormsg = '';
				this.user_or_group_for_send_msg = null;
				this.user_or_group_for_forward_msg = null;
			}
		}
	},
	beforeMount: function () {
        window.addEventListener('beforeunload', (e) => {
			try {
            	// sessionStorage.setItem('user_or_group_for_send_msg',  JSON.stringify(this.user_or_group_for_send_msg));
            	// sessionStorage.setItem('user_or_group_for_forward_msg',  JSON.stringify(this.user_or_group_for_forward_msg));
			} catch {
        		this.$emit('reloginNeeded');
			}
        });

        try {
            // if (sessionStorage.getItem('user_or_group_for_send_msg') != null) this.user_or_group_for_send_msg = JSON.parse(sessionStorage.getItem('user_or_group_for_send_msg'));
            // if (sessionStorage.getItem('user_or_group_for_forward_msg') != null) this.user_or_group_for_forward_msg = JSON.parse(sessionStorage.getItem('user_or_group_for_forward_msg'));
        } catch {
        	this.$emit('reloginNeeded');
		}
    },
	mounted: function() {
		this.user_or_group_for_send_msg = null;
		this.user_or_group_for_forward_msg = null;
	}
}
</script>

<template>
	<div v-if="session_token != 0">
		<div class="message-container">
			<div v-if="selected_conversation != null" style="position: relative; top: 0.7em;">
				<input id="messageInput" v-model="message" placeholder="Message" style="width: 12em;">
                <button @click="doSendMessage">Send</button>
                <span class="label-flat">or</span>
                <label for="photoSender" class="label-button">Send photo</label>
				<input id="photoSender" ref="photoSender" type="file" accept="image/*" hidden="true" @click="onImageUploaderClick" @change="doSendPhoto">
			</div>
			<br>
			<span v-if="selected_message != null && selected_message.from_user_id!=user.id" style="position: relative; top:0.4em">
				<span>&nbsp;&nbsp;&nbsp;</span>
				<button class="emoticon-button" @click="onReactionClick" value="&#x1F600;">&#x1F600;</button>
				<span>&nbsp;</span>
				<button class="emoticon-button" @click="onReactionClick" value="&#x1F602;">&#x1F602;</button>
				<span>&nbsp;</span>
				<button class="emoticon-button" @click="onReactionClick" value="&#x1F62F;">&#x1F62F;</button>
				<span>&nbsp;</span>
				<button class="emoticon-button" @click="onReactionClick" value="&#x1F621;">&#x1F621;</button>
				<span>&nbsp;</span>
				<button class="emoticon-button" @click="onReactionClick" value="&#x1F44D;">&#x1F44D;</button>
				<button @click="onDeleteReaction">Uncomment</button>
			</span>
			<span v-if="selected_message != null && selected_message.from_user_id==user.id" style="position: relative; top: 0.4em;">
				<button @click="onDeleteMessage">Delete msg</button>
			</span>
			<span v-if="selected_message != null" style="position: relative; top: 0.4em;">
				<span class="label-flat">Forward to</span>
				<select id="userOrGroupForForwardMsgSelect" style="position:relative; height: 1.3em; width: 11em;" v-model="user_or_group_for_forward_msg" :selected="0">
					<option style="color:gray" disabled="true" :key="0" :value="null">select user or group</option>
					<option v-for="u in all_users_and_my_groups" :key="(u.is_group ? 'G' : 'U') + u.user_id_or_group_id" :value="u">{{ u.is_group ? u.group_name + " (G)" : u.user_name }}</option>
				</select>
				<button @click="doForwardMessage">Send</button>
			</span>
		</div>
		<ErrorMsg :errormsg="errormsg" @errorDismissed="onErrorDismissed"></ErrorMsg>
	</div>
</template>
