<template>
    <el-table :data="sortedMembers" stripe border>
      <el-table-column prop="rank" label="排名" width="80" align="center" />
      <el-table-column prop="Username" label="用户名" />
      <el-table-column prop="CfRating" label="Codeforces分数" width="150" align="center" />
    </el-table>
  </template>
  
  <script>
  export default {
    name: 'CodeforcesRank',
    props: {
      members: {
        type: Array,
        required: true,
      },
    },
    computed: {
      sortedMembers() {
        const sorted = [...this.members]
          .map(item => ({
            ...item,
            CfRating: Number(item.CfRating) || 0,
          }))
          .sort((a, b) => b.CfRating - a.CfRating)
          .map((item, index) => ({
            ...item,
            rank: index + 1,
          }))
        return sorted
      },
    },
  }
  </script>
  