<template>
    <div class="ranking-page">
        <div class="header">
            <h1 class="header-title">📈 排行榜中心</h1>
            <p class="header-subtitle">追踪队员们在各大平台的精彩表现！</p>
        </div>

        <div class="button-group">
            <el-button v-for="item in tabs" :key="item.value" :type="activeTab === item.value ? 'primary' : 'default'"
                @click="activeTab = item.value" class="platform-btn">
                {{ item.label }}
                <span class="hover-bubble"></span>
            </el-button>
        </div>

        <transition name="fade-slide" mode="out-in">
            <div class="rank-sections" :key="activeTab">
                <NowcoderRank v-if="activeTab === 'nowcoder'" :members="teamMembers" />
                <CodeforcesRank v-if="activeTab === 'codeforces'" :members="teamMembers" />
                <AtcoderRank v-if="activeTab === 'atcoder'" :members="teamMembers" />
            </div>
        </transition>
    </div>
</template>

<script>
import NowcoderRank from '@/components/rank/NowcoderRank.vue'
import CodeforcesRank from '@/components/rank/CodeforcesRank.vue'
import AtcoderRank from '@/components/rank/AtcoderRank.vue'
import request from '@/utils/request'

export default {
    name: 'RankingPage',
    components: {
        NowcoderRank,
        CodeforcesRank,
        AtcoderRank,
    },
    data() {
        return {
            activeTab: 'nowcoder',
            tabs: [
                { label: 'Nowcoder 排名', value: 'nowcoder' },
                { label: 'Codeforces 排名', value: 'codeforces' },
                { label: 'AtCoder 排名', value: 'atcoder' },
            ],
            teamMembers: [],
        }
    },
    created() {
        this.fetchTeamMembers()
    },
    methods: {
        async fetchTeamMembers() {
            try {
                const res = await request.get('/admin/team-members')
                this.teamMembers = res.data.data || []
            } catch (error) {
                this.$message.error('获取团队成员数据失败')
            }
        },
    },
}
</script>

<style scoped>
.ranking-page {
    padding: 2rem 1.5rem;
    max-width: 1200px;
    margin: 0 auto;
}

.header {
    text-align: center;
    margin-bottom: 2.5rem;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 3rem 1rem;
    border-radius: 1rem;
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
    color: white;
}

.header-title {
    font-size: 3rem;
    font-weight: 800;
    margin-bottom: 1rem;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-subtitle {
    font-size: 1.25rem;
    opacity: 0.95;
    font-weight: 300;
}

.button-group {
    text-align: center;
    margin-bottom: 3rem;
    position: relative;
}

.platform-btn {
    position: relative;
    margin: 0 1rem;
    padding: 0.8rem 2rem;
    border-radius: 50px;
    transition: all 0.3s ease;
    overflow: hidden;
}

.platform-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

.hover-bubble {
    position: absolute;
    background: rgba(255, 255, 255, 0.2);
    width: 100px;
    height: 100px;
    border-radius: 50%;
    transform: translate(-50%, -50%) scale(0);
    transition: transform 0.4s ease;
}

.platform-btn:hover .hover-bubble {
    transform: translate(-50%, -50%) scale(1.5);
}

.rank-sections {
    background: #ffffff;
    border-radius: 1.5rem;
    padding: 2rem;
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
    transition: box-shadow 0.3s ease;
}

.rank-sections:hover {
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
}

/* 切换动画 */
.fade-slide-enter-active,
.fade-slide-leave-active {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from {
    opacity: 0;
    transform: translateY(10px);
}

.fade-slide-leave-to {
    opacity: 0;
    transform: translateY(-10px);
}

@media (max-width: 768px) {
    .header-title {
        font-size: 2rem;
    }

    .button-group {
        display: grid;
        gap: 1rem;
    }

    .platform-btn {
        margin: 0;
        width: 100%;
    }
}
</style>