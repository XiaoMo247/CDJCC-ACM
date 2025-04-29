<template>
    <el-table :data="sortedMembers" stripe border>
        <el-table-column prop="rank" label="排名" width="80" align="center" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="nc_id" label="牛客ID" />
        <el-table-column label="牛客分数" width="150" align="center">
            <template #default="scope">
                <span :class="getRatingClass(scope.row.nc_rating)">{{ scope.row.nc_rating || '-' }}</span>
            </template>
        </el-table-column>
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
                    username: item.Username,
                    nc_id: item.NcID || '-',
                    nc_rating: item.NcRating ? Number(item.NcRating) : 0
                }))
                .filter(item => item.nc_id !== '-')
                .sort((a, b) => b.nc_rating - a.nc_rating)
                .map((item, index) => ({
                    ...item,
                    rank: index + 1,
                }))
            return sorted
        },
    },
    methods: {
        getRatingClass(rating) {
            if (!rating) return 'rating-none'
            rating = Number(rating)
            if (rating >= 2400) return 'rating-red'
            if (rating >= 2000) return 'rating-orange'
            if (rating >= 1600) return 'rating-purple'
            if (rating >= 1200) return 'rating-blue'
            if (rating >= 800) return 'rating-green'
            return 'rating-gray'
        }
    }
}
</script>

<style scoped>
.rating-none {
    color: #909399;
}
.rating-gray {
    color: #909399;
}
.rating-green {
    color: #67C23A;
}
.rating-blue {
    color: #409EFF;
}
.rating-purple {
    color: #B620E0;
}
.rating-orange {
    color: #E6A23C;
}
.rating-red {
    color: #F56C6C;
}
</style>
  