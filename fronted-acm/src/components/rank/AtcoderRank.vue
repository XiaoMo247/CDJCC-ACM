<template>
    <el-table :data="sortedMembers" stripe border>
        <el-table-column prop="rank" label="排名" width="80" align="center" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="at_name" label="AtCoder用户名" />
        <el-table-column label="AtCoder分数" width="150" align="center">
            <template #default="scope">
                <span :class="getRatingClass(scope.row.at_rating)">{{ scope.row.at_rating || '-' }}</span>
            </template>
        </el-table-column>
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
                    username: item.Username,
                    at_name: item.AtName || '-',
                    at_rating: item.AtRating ? Number(item.AtRating) : 0
                }))
                .filter(item => item.at_name !== '-')
                .sort((a, b) => b.at_rating - a.at_rating)
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
            if (rating >= 2800) return 'rating-red'
            if (rating >= 2400) return 'rating-orange'
            if (rating >= 2000) return 'rating-yellow'
            if (rating >= 1600) return 'rating-blue'
            if (rating >= 1200) return 'rating-cyan'
            if (rating >= 800) return 'rating-green'
            if (rating >= 400) return 'rating-brown'
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
    color: #808080;
}
.rating-brown {
    color: #804000;
}
.rating-green {
    color: #008000;
}
.rating-cyan {
    color: #00C0C0;
}
.rating-blue {
    color: #0000FF;
}
.rating-yellow {
    color: #C0C000;
}
.rating-orange {
    color: #FF8000;
}
.rating-red {
    color: #FF0000;
}
</style>
  