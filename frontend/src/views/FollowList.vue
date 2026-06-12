<script setup>import { ref, inject, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { getFriendList, getFriendRequests, getSentFriendRequests, acceptFriendRequest, rejectFriendRequest, deleteFriend } from '../api/follow';
const router = useRouter();
const route = useRoute();
const user = inject('user');
const friends = ref([]);
const requests = ref([]);
const sentRequests = ref([]);
const activeTab = ref('friends');
const tabs = [
 { title: '好友', value: 'friends' },
 { title: '好友请求', value: 'requests' },
 { title: '已发送', value: 'sent' }
];
const loadFriends = async () => {
 try {
 const response = await getFriendList();
 friends.value = response.data.friends || [];
 }
 catch (error) {
 console.error('加载好友失败:', error);
 }
};
const loadRequests = async () => {
 try {
 const response = await getFriendRequests();
 requests.value = response.data.requests || [];
 }
 catch (error) {
 console.error('加载好友请求失败:', error);
 }
};
const loadSentRequests = async () => {
 try {
 const response = await getSentFriendRequests();
 sentRequests.value = response.data.sent_requests || [];
 }
 catch (error) {
 console.error('加载已发送请求失败:', error);
 }
};
const handleAccept = async (requestId) => {
 try {
 await acceptFriendRequest(requestId);
 loadRequests();
 loadFriends();
 }
 catch (error) {
 console.error('同意好友请求失败:', error);
 }
};
const handleReject = async (requestId) => {
 try {
 await rejectFriendRequest(requestId);
 loadRequests();
 }
 catch (error) {
 console.error('拒绝好友请求失败:', error);
 }
};
const handleDeleteFriend = async (friendId) => {
 if (!confirm('确定要删除这个好友吗？'))
 return;
 try {
 await deleteFriend(friendId);
 loadFriends();
 }
 catch (error) {
 console.error('删除好友失败:', error);
 }
};
const switchTab = (tab) => {
 activeTab.value = tab;
 if (tab === 'friends')
 loadFriends();
 else if (tab === 'requests')
 loadRequests();
 else if (tab === 'sent')
 loadSentRequests();
};
onMounted(() => {
 if (!user.value) {
 router.push('/login');
 return;
 }
 loadFriends();
});
</script>

<template>
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.push('/profile')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>好友管理</v-toolbar-title>
    </v-app-bar>
    
    <!-- Tab切换 -->
    <v-tabs v-model="activeTab" @update:model-value="switchTab">
      <v-tab v-for="tab in tabs" :key="tab.value" :value="tab.value">
        {{ tab.title }}
        <span v-if="tab.value === 'requests' && requests.length > 0" class="ml-2">
          <span class="badge badge-pill badge-danger">{{ requests.length }}</span>
        </span>
      </v-tab>
    </v-tabs>
    
    <v-container class="py-6">
      <!-- 好友列表 -->
      <template v-if="activeTab === 'friends'">
        <v-card v-if="friends.length > 0">
          <v-list>
            <v-list-item v-for="friend in friends" :key="friend.id">
              <v-list-item-avatar>
                <v-icon color="primary">mdi-account</v-icon>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title>{{ friend.display_name || friend.username }}</v-list-item-title>
                <v-list-item-subtitle>@{{ friend.username }}</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-actions>
                <v-btn text color="error" @click="handleDeleteFriend(friend.id)">删除好友</v-btn>
              </v-list-item-actions>
            </v-list-item>
          </v-list>
        </v-card>
        
        <v-card v-else class="text-center py-12">
          <v-icon size="64" color="grey">mdi-users</v-icon>
          <p class="mt-4 text-grey">暂无好友</p>
        </v-card>
      </template>
      
      <!-- 好友请求 -->
      <template v-if="activeTab === 'requests'">
        <v-card v-if="requests.length > 0">
          <v-list>
            <v-list-item v-for="request in requests" :key="request.id">
              <v-list-item-avatar>
                <v-icon color="primary">mdi-account</v-icon>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title>{{ request.sender.display_name || request.sender.username }}</v-list-item-title>
                <v-list-item-subtitle v-if="request.message">{{ request.message }}</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-actions>
                <v-btn text color="success" @click="handleAccept(request.id)">同意</v-btn>
                <v-btn text color="error" @click="handleReject(request.id)">拒绝</v-btn>
              </v-list-item-actions>
            </v-list-item>
          </v-list>
        </v-card>
        
        <v-card v-else class="text-center py-12">
          <v-icon size="64" color="grey">mdi-user-plus</v-icon>
          <p class="mt-4 text-grey">暂无好友请求</p>
        </v-card>
      </template>
      
      <!-- 已发送请求 -->
      <template v-if="activeTab === 'sent'">
        <v-card v-if="sentRequests.length > 0">
          <v-list>
            <v-list-item v-for="request in sentRequests" :key="request.id">
              <v-list-item-avatar>
                <v-icon color="primary">mdi-account</v-icon>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title>{{ request.receiver.display_name || request.receiver.username }}</v-list-item-title>
                <v-list-item-subtitle>等待对方确认</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-actions>
                <v-chip color="warning" text-color="white" size="small">待处理</v-chip>
              </v-list-item-actions>
            </v-list-item>
          </v-list>
        </v-card>
        
        <v-card v-else class="text-center py-12">
          <v-icon size="64" color="grey">mdi-send</v-icon>
          <p class="mt-4 text-grey">暂无已发送的好友请求</p>
        </v-card>
      </template>
    </v-container>
  </v-app>
</template>