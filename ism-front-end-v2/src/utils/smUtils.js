
import { sm2, sm3, sm4 } from 'sm-crypto'

// SM2加密解密（Base64格式）
const sm2PublicKey = '公钥Base64字符串'
const sm2PrivateKey = '私钥Base64字符串'

export const sm2EncryptBase64 = (data) => {
    const encrypted = sm2.doEncrypt(data, sm2PublicKey)
    return Buffer.from(encrypted, 'hex').toString('base64')
}

export const sm2DecryptBase64 = (base64Str) => {
    const encrypted = Buffer.from(base64Str, 'base64').toString('hex')
    return sm2.doDecrypt(encrypted, sm2PrivateKey)
}

// SM4加密解密（Base64格式）
const sm4Key = 'c0a419015e6abb6981662ad51d01b780'

export const sm4EncryptBase64 = (data) => {
    const encrypted = sm4.encrypt(data, sm4Key)
    return Buffer.from(encrypted, 'hex').toString('base64')
}

export const sm4DecryptBase64 = (base64Str) => {
    const encrypted = Buffer.from(base64Str, 'base64').toString('hex')
    return sm4.decrypt(encrypted, sm4Key)
}

// SM3哈希（Base64输出）
export const sm3HashBase64 = (data) => {
    const hash = sm3(data)
    return Buffer.from(hash, 'hex').toString('base64')
}
