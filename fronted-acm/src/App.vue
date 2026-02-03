<template>
  <div id="app">
    <a href="#main-content" class="skip-link">跳到主内容</a>
    <NavBar />
    <main id="main-content" class="main-content" role="main" aria-label="主要内容">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" :key="$route.fullPath" />
        </transition>
      </router-view>
    </main>
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
  height: 100vh;
  overflow: hidden; /* 禁止App层级滚动 */
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

/* 跳过导航链接（键盘/读屏用户） */
.skip-link {
  position: absolute;
  top: -40px;
  left: 0;
  background: var(--color-primary, #3498db);
  color: var(--color-white, #fff);
  padding: 8px 16px;
  text-decoration: none;
  z-index: 1100;
}

.skip-link:focus {
  top: 0;
}

/* 主体内容区，避免被导航栏遮住 */
.main-content {
  flex: 1;
  padding-top: 72px; /* 与 NavBar 高度一致 */
  padding-bottom: 40px; /* 留出 Footer 的空间 */
  box-sizing: border-box;
  overflow: hidden; /* 禁止main-content层级滚动，让内部组件自己处理 */
}

</style>
