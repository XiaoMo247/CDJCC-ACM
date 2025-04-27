<template>
    <el-table :data="sortedMembers" stripe border>
      <el-table-column prop="rank" label="排名" width="80" align="center" />
      <el-table-column prop="Username" label="用户名" />
      <el-table-column prop="NcRating" label="牛客分数" width="150" align="center" />
    </el-table>
  </template>
  
  <script>
  export default {
    name: 'NowcoderRank',
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
            NcRating: Number(item.NcRating) || 0, // 👈 强制转成数字
          }))
          .sort((a, b) => b.NcRating - a.NcRating)
          .map((item, index) => ({
            ...item,
            rank: index + 1,
          }))
        return sorted
      },
    },
  }
  </script>
  