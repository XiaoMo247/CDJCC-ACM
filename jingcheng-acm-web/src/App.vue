<template>
  <!-- 主应用容器 -->
  <div id="app" class="min-h-screen flex flex-col bg-gray-50">
    <!-- 导航栏 -->
    <NavBar class="shadow-md" />
    
    <!-- 主要内容区域 -->
    <main class="flex-grow container mx-auto px-4 py-8">
      <!-- 页面过渡效果 -->
      <transition 
        name="page-fade"
        mode="out-in"
      >
        <router-view class="bg-white rounded-lg shadow-sm p-6" />
      </transition>
    </main>
    
    <!-- 页脚 -->
    <Footer class="bg-gray-800 text-white py-6" />
  </div>
</template>

<script setup>
// 组件导入
import NavBar from '@/components/NavBar.vue'
import Footer from '@/components/Footer.vue'
</script>

<style>
/* 全局页面切换动画 - 更平滑的渐变效果 */
.page-fade-enter-active {
  transition: all 0.3s ease-out;
}

.page-fade-leave-active {
  transition: all 0.2s ease-in;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

/* 现代加载动画 */
@keyframes rotate {
  0% {
    transform: rotate(0deg);
    border-top-color: #3498db;
  }
  50% {
    border-top-color: #9b59b6;
  }
  100% {
    transform: rotate(360deg);
    border-top-color: #3498db;
  }
}

.loading-spinner {
  display: inline-block;
  animation: rotate 1.5s cubic-bezier(0.68, -0.55, 0.27, 1.55) infinite;
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-top: 4px solid;
  border-radius: 50%;
  width: 48px;
  height: 48px;
  margin: 2rem auto;
  position: relative;
}

.loading-spinner:after {
  content: "";
  position: absolute;
  top: -4px;
  left: -4px;
  right: -4px;
  bottom: -4px;
  border-radius: 50%;
  border: 4px solid transparent;
  border-bottom-color: rgba(52, 152, 219, 0.3);
}

/* 基础样式重置 */
#app {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>