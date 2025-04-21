<template>
  <div id="app">
    <NavBar />
    <div class="main-content">
      <router-view />
    </div>
    <Footer />
  </div>
</template>

<script>
import NavBar from './components/NavBar.vue'
import Footer from './components/Footer.vue'

export default {
  name: 'App',
  components: {
    NavBar,
    Footer
  }
}
</script>

<style>
#app {
  padding-top: 20px; /* 留出顶部导航栏高度 */
  padding-bottom: 20px; /* 留出底部 Footer 高度 */
  min-height: 100vh;
  box-sizing: border-box;
}

.main-content {
  padding: 20px;
}
</style>
