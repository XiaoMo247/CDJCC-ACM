const express = require('express');
const router = express.Router();
const multer = require('multer');
const path = require('path');
const fs = require('fs');
const { verifyToken } = require('../middleware/auth');
const axios = require('axios');

// 配置multer存储
const storage = multer.diskStorage({
    destination: function (req, file, cb) {
        const uploadDir = 'public/uploads/avatars';
        // 确保上传目录存在
        if (!fs.existsSync(uploadDir)) {
            fs.mkdirSync(uploadDir, { recursive: true });
        }
        cb(null, uploadDir);
    },
    filename: function (req, file, cb) {
        // 生成唯一文件名：时间戳 + 随机数 + 原扩展名
        const uniqueSuffix = Date.now() + '-' + Math.round(Math.random() * 1E9);
        cb(null, uniqueSuffix + path.extname(file.originalname));
    }
});

// 文件过滤器
const fileFilter = (req, file, cb) => {
    // 只允许图片文件
    if (file.mimetype.startsWith('image/')) {
        cb(null, true);
    } else {
        cb(new Error('只能上传图片文件！'), false);
    }
};

const upload = multer({
    storage: storage,
    fileFilter: fileFilter,
    limits: {
        fileSize: 5 * 1024 * 1024 // 限制5MB
    }
});

// 头像上传接口
router.post('/upload-avatar', verifyToken, upload.single('avatar'), async (req, res) => {
    try {
        if (!req.file) {
            return res.status(400).json({
                code: 400,
                msg: '请选择要上传的图片'
            });
        }

        // 获取上传的文件信息
        const file = req.file;
        
        // 构建图片URL
        const avatarUrl = `/uploads/avatars/${file.filename}`;

        res.json({
            code: 200,
            msg: '头像上传成功',
            data: {
                avatarUrl: avatarUrl
            }
        });
    } catch (error) {
        console.error('头像上传失败:', error);
        res.status(500).json({
            code: 500,
            msg: '头像上传失败'
        });
    }
});

// 获取QQ头像接口
router.get('/qq-avatar/:qq', verifyToken, async (req, res) => {
    try {
        const qq = req.params.qq;
        if (!qq) {
            return res.status(400).json({
                code: 400,
                msg: 'QQ号不能为空'
            });
        }

        // 调用QQ头像API
        const response = await axios.get(`https://q1.qlogo.cn/g?b=qq&nk=${qq}&s=640`, {
            responseType: 'arraybuffer'
        });

        // 将图片转换为base64
        const base64 = Buffer.from(response.data).toString('base64');
        const avatarUrl = `data:image/jpeg;base64,${base64}`;

        res.json({
            code: 200,
            msg: '获取头像成功',
            data: {
                avatarUrl: avatarUrl
            }
        });
    } catch (error) {
        console.error('获取QQ头像失败:', error);
        // 如果获取失败，返回默认头像
        res.json({
            code: 200,
            msg: '使用默认头像',
            data: {
                avatarUrl: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAOEAAADhCAMAAAAJbSJIAAAAllBMVEW/2OvF4vra5+shISFCQkLH4vnA2u7T5fDN5PXF3/O+1+vI5v5SW2IbGRgNBACZrsAHAADN6/+rxNkeHRw1OT02My/W5u7i7/OosrUXFBI4OjumvMyUp7WKmaU7OTdRU1N4hI1kbHIzMzNreIEuLStKUlexy+AqLTDI3euQo7GGlKBMTEsuKCHY6/Z/jJV2gYokHxpDRkrDp6E3AAACnklEQVR4nO3d0VbaQABF0ahVCUHBAFosRUAJVKm0//9z/YJ7H+7KWNt1znsyswN5yaxJqvMCDc50g/Cc9+acl+7AKhzQhjAJoQxhFMIkhDKEUQiTEMoQRiFMQihDGIUwCaEMYRTCJIQyhFEIkz5e+GBywutwvDLCSldvto+ql3FjZpNmhGdmolUsHC2mqokTFqmM8OZCdYsQIUKECBEiRIgQIUKECBEiRIgQIUKECP8n4W6iKyK80zVWWJvcMkL7TdddZdUG2O1Gsk1jJlqll9TN1F5S05UermknN6rF7M5MNBemjFgo7/spQoQIESJEiBAhQoQIESJEiBAhQoQIESL8lMKhzgtvVZO9FV6bSjy5P7gBTYcfS13tJhquPaUN4z0z8eX+Z4TDdEiEPYcQYRDCnkOIMAhhzyFEGISw5xAiDELYcwgRBiHsOYQIg9xZ7cPysEMs1MXCh++6185t1DANTO73vTcLHodY+LaSHdsC+2LSN52d21/RCZ++yOZt/0CECBEiRIgQIUKECBEiRIgQIUKECBEiRPjXhO6Zd9foigjN9zS80Gxx+Pm+lp2WY10XCodfw+w6kRvRvSRrZ/awhO/6asbHrezFvFtsX6dCc1iJ97U147m+8R/lcBfTGUKECBEiRIgQIUKECBEiRIgQIUKECBEiLClcyG8V+GfeZr3DCVd6vMXWCt2Te3fYSX9uYrTs2qjlxmTG21ihW+5xxzXmmyHr41z2y6w/PF+G21TcRGOh+dZ5877S/7atvp8Wz36qYQgRIkSIECFChAgRIkSIECFChAgRIkSIsJxw7YRT2WcT6prXtyfZ75lsf/pwYdrArGm43FJQjCgj1P9gG0KECBEiRIgQIUKECBEiRIgQIUKECBEiRIjQ9AeceaXi/L2vVwAAAABJRU5ErkJggg=='
            }
        });
    }
});

// ... 其他路由保持不变 ...

module.exports = router; 