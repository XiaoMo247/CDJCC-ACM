<template>
    <div class="ranking-page page-ranking">
        <div class="header page-header page-hero">
            <h1 class="header-title">📈 排行榜中心</h1>
            <p class="header-subtitle">追踪队员们在各大平台的精彩表现！</p>
            <el-button type="primary" :loading="loading" @click="refreshData" class="refresh-btn">
                <el-icon>
                    <Refresh />
                </el-icon>
                刷新数据
            </el-button>
        </div>

        <div class="button-group">
            <el-button v-for="item in tabs" :key="item.value" :type="activeTab === item.value ? 'primary' : 'default'"
                @click="activeTab = item.value" class="platform-btn">
                {{ item.label }}
                <span class="hover-bubble"></span>
            </el-button>
        </div>

        <transition name="fade-slide" mode="out-in">
            <div v-if="loading" class="loading-container">
                <el-skeleton :rows="10" animated />
            </div>
            <div v-else-if="error" class="error-container">
                <el-empty :image-size="200" :description="error">
                    <template #extra>
                        <el-button type="primary" @click="refreshData">重试</el-button>
                    </template>
                </el-empty>
            </div>
            <div v-else class="rank-sections" :key="activeTab">
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
import { Refresh } from '@element-plus/icons-vue'

export default {
    name: 'RankingPage',
    components: {
        NowcoderRank,
        CodeforcesRank,
        AtcoderRank,
        Refresh
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
            loading: false,
            error: null,
            refreshInterval: null
        }
    },
    created() {
        this.fetchTeamMembers()
        // 每5分钟自动刷新一次
        this.refreshInterval = setInterval(this.fetchTeamMembers, 300000)
    },
    beforeUnmount() {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval)
        }
    },
    methods: {
        async fetchTeamMembers() {
            this.loading = true
            this.error = null
            try {
                const res = await request.get('/admin/team-members')
                if (res.data && res.data.data) {
                    this.teamMembers = res.data.data
                } else {
                    throw new Error('数据格式错误')
                }
            } catch (error) {
                console.error('获取团队成员数据失败:', error)
                this.error = error.response?.data?.msg || '获取数据失败，请稍后重试'
            } finally {
                this.loading = false
            }
        },
        refreshData() {
            this.fetchTeamMembers()
        }
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
    margin-bottom: 60px;
    background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
    padding: 60px 20px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    color: #78350f;
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
}

.refresh-btn {
    position: absolute;
    right: 2rem;
    top: 2rem;
    border-radius: 50px;
    padding: 0.6rem 1.2rem;
    background: rgba(255, 255, 255, 0.2);
    border: none;
    color: white;
    transition: all 0.3s ease;
}

.refresh-btn:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: scale(1.05);
}

.header-title {
    font-size: 3rem;
    font-weight: 800;
    margin-bottom: 1rem;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    display: inline-flex;
    align-items: center;
}

.header-title::before {
    content: "🥇";
    width: 44px;
    height: 44px;
    border-radius: 12px;
    background: rgba(255, 255, 255, 0.16);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    margin-right: 12px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
}

.header-subtitle {
    font-size: 1.25rem;
    opacity: 0.95;
    font-weight: 300;
}

.button-group {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 1rem;
    margin-bottom: 3rem;
    width: 100%;
    box-sizing: border-box;
}

.platform-btn {
    position: relative;
    padding: 0.8rem 2rem;
    border-radius: 50px;
    transition: all 0.3s ease;
    overflow: hidden;
    box-sizing: border-box;
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

.loading-container,
.error-container {
    background: #ffffff;
    border-radius: 1.5rem;
    padding: 2rem;
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
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

/* ✅ 移动端适配优化 */
@media (max-width: 768px) {
    .header-title {
        font-size: 2rem;
    }

    .header {
        padding: 2rem 1rem 1rem;
    }

    .refresh-btn {
        position: relative;
        margin-top: 1.5rem;
        right: auto;
        top: auto;
        width: 100%;
        max-width: 200px;
        align-self: center;
    }

    .button-group {
        flex-direction: column;
        align-items: center;
        gap: 1rem;
    }

    .platform-btn {
        width: 300px;
        margin: 0 10px;
        text-align: center;
    }
}
</style>
