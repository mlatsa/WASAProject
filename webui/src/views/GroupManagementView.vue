<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    emits: ['groupsUpdated', 'reloginNeeded'],
	props: ['session_token', 'all_users', 'my_groups'],
	data: function() {
		return {
            errormsg: '',
            group_name_to_create: '',
            user_for_new_group: null,
            user_to_add: null,
            group_name_to_growth: null,
            group_to_set: null,
            new_group_name: ''
		}
	},
	methods: {
        async doCreateGroup() {
            try {
                if (this.group_name_to_create == '') {
                    this.errormsg = "Group name not set";
                    return;
                }
                
                if (this.user_for_new_group == null) {
                    this.errormsg = "User with which create a group not set";
                    return;
                }
                
                const res = await this.$axios({
                    method: 'put',
                    url: '/groups/0',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "user_name_to_add": this.user_for_new_group.user_name,
                        "group_name": this.group_name_to_create
                    }
                });

                if (res.status != 201) {
                    this.errormsg = res.data;
                    return;
                }

                this.$emit('groupsUpdated');

                this.group_name_to_create = '';
                this.user_for_new_group = null

                this.errormsg = "";
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        async doAddToGroup() {
            try {
                if (this.user_to_add == null) {
                    this.errormsg = "User to add not set";
                    return;
                }
                
                if (this.group_to_growth == null) {
                    this.errormsg = "Group name not set";
                    return;
                }
                
                const res = await this.$axios({
                    method: 'put',
                    url: '/groups/0',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "user_name_to_add": this.user_to_add.user_name,
                        "group_name": this.group_to_growth.group_name
                    }
                });

                if (res.status != 201) {
                    this.errormsg = res.data;
                    return;
                }

                this.$emit('groupsUpdated');

                this.user_to_add = null
                this.group_to_growth = null;

                this.errormsg = "";
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        async doSetGroupName() {
            try {
                if (this.group_to_set == null) {
                    this.errormsg = "Group name not set";
                    return;
                }

                if (this.new_group_name == '') {
                    this.errormsg = "New user group not set";
                    return;
                }
                
                const res = await this.$axios({
                    method: 'post',
                    url: '/groups/' + this.group_to_set.group_id,
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "new_group_name": this.new_group_name
                    }
                });

                if (res.status != 200) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                this.$emit('groupsUpdated');
                this.new_group_name = '';
                this.group_to_set = null;

                this.errormsg = "";
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        // this code cannot be placed inside doSetMyPhoto because
        // await need to be on the top level of an async method
        async doUploadGroupPhoto(img)
        {
            try {
                const res = await this.$axios({
                    method: 'put',
                    url: '/group/photo',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "group_id": this.group_to_set.group_id,
                        "photo": img
                    }
                });

                if (res.status != 200) {
                    this.errormsg = "Problem while updating the group photo";
                }

                this.$emit('groupsUpdated');
                this.errormsg = "";
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },
		
        async doSetGroupPhoto(e) {
            if (this.group_to_set == null) {
                this.errormsg = "Group name not set";
                return;
            }

            const img = e.target.files[0];
            const reader = new FileReader();
            reader.readAsDataURL(img);
            reader.onload = evn =>{
                this.doUploadGroupPhoto(evn.target.result);
                this.group_to_set.group_photo = evn.target.result;
            };
        },

        async onPhotoGroupUploaderClick() {
            this.$refs.photoGroupUploader.value = ''
        },

        async doLeaveGroup() {
            try {
                if (this.group_to_set == null) {
                    this.errormsg = "Group to leave not set";
                    return;
                }

                const res = await this.$axios({
                    method: 'delete',
                    url: '/groups/' + this.group_to_set.group_id,
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                });

                if (res.status != 200) {
                    this.errormsg = "Problem while leaving the group";
                }

                this.$emit('groupsUpdated');
                this.errormsg = "";
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },
        async doShowMembers() {
            try {
                if (this.group_to_set == null) {
                    this.errormsg = "Group not set";
                    return;
                }

                const res = await this.$axios({
                    method: 'get',
                    url: '/groups/' + this.group_to_set.group_id,
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                });

                if (res.status != 200) {
                    this.errormsg = "Problem while retriving members of the group";
                }
                
                let membersList = '';
                for (var i = 0; i < res.data.length; i++) {
                    membersList += res.data[i]['user_name'] + '\n';
                }
                alert("I membri del gruppo sono: \n" + membersList)

                this.errormsg = "";
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
      		if (this.session_token == 0) {
                this.group_name_to_create = '';
                this.user_for_new_group = null;
                this.user_to_add = null;
                this.group_to_growth = null;
                this.group_to_set = null;
                this.new_group_name = '';
                this.errormsg = '';
            }
		},
        all_users(newValue, oldValue) {
            this.user_for_new_group = null;
            this.user_to_add = null;
        },
        my_groups(newValue, oldValue) {
            this.group_to_growth = null;
            this.group_to_set = null;
        }
    },
    beforeMount: function () {
        window.addEventListener('beforeunload', (e) => {
            try {
                // sessionStorage.setItem('group_name_to_create',  JSON.stringify(this.group_name_to_create));
                // sessionStorage.setItem('user_for_new_group', JSON.stringify(this.user_for_new_group));
                // sessionStorage.setItem('user_to_add', JSON.stringify(this.user_to_add));
                // sessionStorage.setItem('group_to_growth', JSON.stringify(this.group_to_growth));
                // sessionStorage.setItem('group_to_set', JSON.stringify(this.group_to_set));
                // sessionStorage.setItem('new_group_name', JSON.stringify(this.new_group_name));
            } catch {
				this.$emit('reloginNeeded');
			}
        });

        try {
            // if (sessionStorage.getItem('group_name_to_create') != null) this.group_name_to_create = JSON.parse(sessionStorage.getItem('group_name_to_create'));
            // if (sessionStorage.getItem('user_for_new_group') != null) this.user_for_new_group = JSON.parse(sessionStorage.getItem('user_for_new_group'));
            // if (sessionStorage.getItem('user_to_add') != null) this.user_to_add = JSON.parse(sessionStorage.getItem('user_to_add'));
            // if (sessionStorage.getItem('group_to_growth') != null) this.group_to_growth = JSON.parse(sessionStorage.getItem('group_to_growth'));
            // if (sessionStorage.getItem('group_to_set') != null) this.group_to_set = JSON.parse(sessionStorage.getItem('group_to_set'));
            // if (sessionStorage.getItem('new_group_name') != null) this.new_group_name = JSON.parse(sessionStorage.getItem('new_group_name'));
        } catch {
        	this.$emit('reloginNeeded');
		}
    },
    mounted: function() {
        this.user_for_new_group = null;
        this.user_to_add = null;
        this.group_to_growth = null;
        this.group_to_set = null;
        this.errormsg = '';
    }
}
</script>

<template>
    <div v-if="session_token != 0">
        <div class="group-management-container">
            <div style="position:relative; top: 0.7em; float: left;">
                <span class="label-flat">Create group</span>
                <input id="groupNameToCreateInput" v-model="group_name_to_create" placeholder="New group" autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="false">
                <span class="label-flat">with</span>
                <select id="userNameForNewGroupSelect" style="position:relative; width: 10em;" v-model="user_for_new_group" :selected="0">
                    <option style="color:gray" disabled="true" :key="0" :value="null">select user</option>
                    <option v-for="u in all_users" :key="u.group_id" :value="u">{{ u.user_name }}</option>
                </select>
                <button @click="doCreateGroup">Apply</button>
                <br>
                <p style="border-bottom: 1em;"></p>
                <span class="label-flat">Add to group</span>
                <select id="groupNameToGrowthSelect" style="position:relative; width: 9em;" v-model="group_to_growth" :selected="0">
                    <option style="color:gray" disabled="true" :key="0" :value="null">select group</option>
                    <option v-for="u in my_groups" :key="u.group_id" :value="u">{{ u.group_name }}</option>
                </select>
                <select id="userNameToAddSelect" style="position:relative; width: 10em;" v-model="user_to_add" :selected="0">
                    <option style="color:gray" disabled="true" :key="0" :value="null">select user</option>
                    <option v-for="u in all_users" :key="u.user_id" :value="u">{{ u.user_name }}</option>
                </select>
                <button @click="doAddToGroup">Apply</button>
                <br>
                <p style="border-bottom: 1em;"></p>
                <span class="label-flat">For group</span>
                <select id="groupNameToSetSelect" style="position:relative; width: 9.5em;" v-model="group_to_set" :selected="0">
                    <option style="color:gray" disabled="true" :key="0" :value="null">select group</option>
                    <option v-for="g in my_groups" :key="g.group_id" :value="g">{{ g.group_name }}</option>
                </select>
                <input id="newGroupNameInput" v-model="new_group_name" placeholder="New name" autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="false">
			    <button @click="doSetGroupName">Apply</button>
                <span class="label-flat">or</span>
                <label for="photoGroupUploader" class="label-button">Set photo</label>
                <input id="photoGroupUploader" type="file" accept="image/*" hidden="true" ref="photoGroupUploader" @click="onPhotoGroupUploaderClick" @change="doSetGroupPhoto">
                <button @click="doLeaveGroup">Leave group</button>
                <button @click="doShowMembers">Show members</button>
            </div>
        </div>
		<ErrorMsg :errormsg="errormsg" @errorDismissed="onErrorDismissed"></ErrorMsg>
    </div>
</template>
