<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	emits: ['messageSent'],
    props: ['session_token', 'user', 'all_users'],
    data: function() {
		return {
			errormsg: '',
			user_to_send_msg: null,
			message: ''
		}
	},
	methods: {
		async doSendMessage() {
            try {
                if (this.user_to_send_msg == null) {
                    this.errormsg = "Receiver of message not set";
                    return;
                }

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
                        "to_user_name_or_group_name": this.user_to_send_msg.user_name,
                        "is_group": false,
                        "message": this.message
                    }
                });

                if (res.status != 201) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                this.$emit('messageSent', this.user_to_send_msg.user_name);
                this.message = '';

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
				this.user_to_send_msg = null;
			}
		},
	},
	beforeMount: function () {
        window.addEventListener('beforeunload', (e) => {
			try {
                sessionStorage.setItem('message',  JSON.stringify(this.message));
			} catch {
        		this.$emit('reloginNeeded');
			}
        });

        try {
            if (sessionStorage.getItem('message') != null) this.message = JSON.parse(sessionStorage.getItem('message'));
        } catch {
        	this.$emit('reloginNeeded');
		}
    },
	mounted: function() {
	}
}
</script>

<template>
    <div v-if="session_token != 0">
		<div class="new-conversation-container" v-if="session_token != 0">
            <div style="position:relative; top: 0.7em; float: left;">
                <span class="label-flat">New conversation with</span>
                <select id="userToSendMsgSelect" style="position:relative; width: 10em;" v-model="user_to_send_msg" :selected="0">
                    <option style="color:gray" disabled="true" :key="0" :value="null">select user</option>
                    <option v-for="u in all_users" :key="u.user_id" :value="u">{{ u.user_name }}</option>
                </select>
                <input id="messageInput" v-model="message" placeholder="Hello Message">
                <button @click="doSendMessage">Send</button>
            </div>
        </div>
    </div>
    <ErrorMsg :errormsg="errormsg" @errorDismissed="onErrorDismissed"></ErrorMsg>
</template>
