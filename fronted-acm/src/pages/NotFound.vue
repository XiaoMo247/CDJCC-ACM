<template>
  <div class="not-found-page">
    <div class="error-container">
      <div class="error-code">404</div>
      <p class="error-message">抱歉，您访问的页面不存在</p>
      <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'NotFound'
}
</script>

<style scoped>
.not-found-page {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-6, 2rem);
}

.error-container {
  text-align: center;
}

.error-code {
  font-size: 120px;
  font-weight: 800;
  line-height: 1;
  color: var(--color-primary, #3498db);
  margin-bottom: var(--spacing-4, 1rem);
}

.error-message {
  font-size: 20px;
  color: var(--color-gray-700, #334155);
  margin-bottom: var(--spacing-5, 1.5rem);
}
</style>
