<template>
    <el-table :data="sortedMembers" stripe border>
      <el-table-column prop="rank" label="排名" width="80" align="center" />
      <el-table-column prop="Username" label="用户名" />
      <el-table-column prop="AtRating" label="AtCoder分数" width="150" align="center" />
    </el-table>
  </template>
  
  <script>
  export default {
    name: 'AtcoderRank',
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
            AtRating: Number(item.AtRating) || 0,
          }))
          .sort((a, b) => b.AtRating - a.AtRating)
          .map((item, index) => ({
            ...item,
            rank: index + 1,
          }))
        return sorted
      },
    },
  }
  </script>
  