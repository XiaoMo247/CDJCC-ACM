<template>
  <div class="footer">
    <p>&copy; 2025 成都锦城学院 ACM 集训队. All Rights Reserved.</p>
  </div>
</template>

<script>
export default {
  name: 'Footer'
}
</script>

<style scoped>
.footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: linear-gradient(90deg, #f1f5f9, #e2e8f0, #cbd5e1);
  color: #1e293b;
  text-align: center;
  line-height: 60px;
  z-index: 1000;
  box-shadow: 0 -4px 6px -1px rgba(0, 0, 0, 0.1);
}

.footer p {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
}
</style>
