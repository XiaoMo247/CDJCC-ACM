<template>
  <div class="contests-container">
    <h1>比赛信息</h1>
    <el-row :gutter="20" justify="center">
      <el-col :span="8" v-for="(contest, index) in contests" :key="index">
        <el-card class="contest-card" @click="goToContest(contest.id)">
          <h2>{{ contest.title }}</h2>
          <p class="date">{{ contest.date }}</p>
          <p>{{ contest.description }}</p>
          <el-button type="primary" class="contest-btn">报名参加</el-button>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
const contests = [
  {
    id: 1,
    title: '校内选拔赛',
    date: '2025-05-01',
    description: '本次比赛为校内选拔赛，旨在选拔优异的算法人才。',
  },
  {
    id: 2,
    title: 'ACM 省赛',
    date: '2025-06-15',
    description: '参与省级ACM编程竞赛，挑战全国顶尖选手。',
  },
]

const goToContest = (id) => {
  console.log('报名比赛', id)
}
</script>

<style scoped>
.contests-container {
  padding: 40px 20px;
}

.contest-card {
  padding: 20px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.contest-card:hover {
  transform: translateY(-5px);
}

.contest-btn {
  margin-top: 10px;
}
</style>
