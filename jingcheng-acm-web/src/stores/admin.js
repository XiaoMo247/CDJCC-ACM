import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAdminStore = defineStore('admin', () => {
  const isLoggedIn = ref(localStorage.getItem('adminToken') !== null)

  const login = (token) => {
    localStorage.setItem('adminToken', token)
    isLoggedIn.value = true
  }

  const logout = () => {
    localStorage.removeItem('adminToken')
    isLoggedIn.value = false
  }

  return {
    isLoggedIn,
    login,
    logout
  }
})