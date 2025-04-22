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
  min-height: 100vh;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

/* 主体内容区，避免被导航栏遮住 */
.main-content {
  flex: 1;
  padding-top: 72px; /* 与 NavBar 高度一致 */
  padding-bottom: 40px; /* 留出 Footer 的空间 */
  box-sizing: border-box;
}
</style>
