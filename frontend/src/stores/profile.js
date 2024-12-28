import { defineStore } from 'pinia';

export const useProfileStore = defineStore('profile', {
  state: () => ({
    userProfile: {
      username: '',
      email: '',
      role: 0
    }
  }),
  actions: {
    setUserProfile(profile) {
      this.userProfile = profile;
    },
    clearUserProfile() {
      this.userProfile = {
        username: '',
        email: '',
        role: 0
      };
    }
  },
  getters: {
    isAuthenticated: (state) => !!state.userProfile.email,
    isAdmin: (state) => state.userProfile.role == 2,
    isSuperAdmin: (state) => state.userProfile.role === 3,
    username: (state) => state.userProfile.username,
  },
  persist: true  // 开启持久化
});