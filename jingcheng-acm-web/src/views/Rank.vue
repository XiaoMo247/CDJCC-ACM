<template>
    <div class="rank-container">
      <h1>排行榜</h1>
      <el-table :data="rankList" style="width: 100%">
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="score" label="积分" />
        <el-table-column prop="rank" label="排名" />
      </el-table>
    </div>
  </template>
  
  <script setup>
  const rankList = [
    { rank: 1, name: '张三', score: 120 },
    { rank: 2, name: '李四', score: 110 },
    { rank: 3, name: '王五', score: 100 },
  ]
  </script>
  
  <style scoped>
  .rank-container {
    padding: 40px 20px;
  }
  </style>
  