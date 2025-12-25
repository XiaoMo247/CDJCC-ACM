import acmIcon from '../../assets/acm-logo.jpg'

import photo1 from '../../assets/ACMALL.jpg'
import photo2 from '../../assets/ACMALL.jpg'
import photo3 from '../../assets/ACMALL.jpg'

import member1 from '../../assets/user.png'
import member2 from '../../assets/user.png'
import member3 from '../../assets/user.png'
import member4 from '../../assets/user.png'
import member5 from '../../assets/user.png'
import member6 from '../../assets/user.png'
import member7 from '../../assets/user.png'
import member8 from '../../assets/user.png'
import member9 from '../../assets/user.png'

// 荣誉数据
export const honors = [
  {
    id: 1,
    name: "张三",
    year: 2023,
    competition: "ICPC",
    level: "国一",
    grade: "大三",
    major: "计算机科学",
    photo: member1
  },
  {
    id: 2,
    name: "李四",
    year: 2022,
    competition: "CCPC",
    level: "国二",
    grade: "大四",
    major: "软件工程",
    photo: member2
  },
  {
    id: 3,
    name: "王五",
    year: 2023,
    competition: "蓝桥杯",
    level: "省一",
    grade: "大二",
    major: "人工智能",
    photo: member3
  }
]

// 团队信息
export const teamInfo = {
    name: '锦城ACM集训队',
    college: '成都锦城学院计算机与软件学院',
    established: '2020年',
    avatar: acmIcon,
    description: '本团队致力于算法学习与ACM竞赛训练，培养团队协作精神与编程能力。',
    awards: 28,
    mission: '培养优秀的程序设计人才，在竞赛中展现锦城风采。',
    features: [
        {
            icon: 'User',
            title: '专业指导',
            desc: '由资深算法竞赛指导老师带队'
        },
        {
            icon: 'Clock',
            title: '系统训练',
            desc: '每周固定训练时间，循序渐进提升能力'
        },
        {
            icon: 'Trophy',
            title: '竞赛实战',
            desc: '定期参加各类高水平程序设计竞赛'
        },
        {
            icon: 'Promotion',
            title: '成长空间',
            desc: '提供广阔的技术发展和就业机会'
        }
    ]
}

// 照片数据
export const photos = [
    {
        url: photo1,
        title: '2023年ICPC区域赛',
        description: '团队参加2023年ICPC亚洲区域赛现场',
        date: '2023-11-15',
        location: '成都'
    },
    {
        url: photo2,
        title: '团队训练',
        description: '日常团队训练场景',
        date: '2023-09-10',
        location: '学校实验室'
    },
    {
        url: photo3,
        title: '获奖合影',
        description: '团队在省级比赛中获得一等奖',
        date: '2023-05-20',
        location: '重庆'
    }
]