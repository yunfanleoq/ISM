/**
 * IndexedDB 工具函数
 * 用于保存/读取大型数据（如3D模型）
 */

const DB_NAME = 'ISM3DEditorDB'
const DB_VERSION = 1
const STORE_NAME = 'models'

/**
 * 打开数据库
 */
function openDB() {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION)
    
    request.onerror = function() {
      reject(new Error('打开数据库失败'))
    }
    
    request.onsuccess = function() {
      resolve(request.result)
    }
    
    request.onupgradeneeded = function(e) {
      const db = e.target.result
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        db.createObjectStore(STORE_NAME, { keyPath: 'id' })
      }
    }
  })
}

/**
 * 保存模型数据到 IndexedDB
 * @param {string} id - 模型ID
 * @param {Object} modelData - 模型数据
 * @param {string} fileName - 文件名
 * @returns {Promise}
 */
export function saveModelToDB(id, modelData, fileName) {
  return openDB().then(db => {
    return new Promise((resolve, reject) => {
      const transaction = db.transaction([STORE_NAME], 'readwrite')
      const store = transaction.objectStore(STORE_NAME)
      
      const data = {
        id: id,
        fileName: fileName,
        data: modelData,
        timestamp: Date.now()
      }
      
      const request = store.put(data)
      
      request.onsuccess = function() {
        resolve()
      }
      
      request.onerror = function() {
        reject(new Error('保存模型失败'))
      }
    })
  })
}

/**
 * 从 IndexedDB 读取模型数据
 * @param {string} id - 模型ID
 * @returns {Promise<Object>}
 */
export function getModelFromDB(id) {
  return openDB().then(db => {
    return new Promise((resolve, reject) => {
      const transaction = db.transaction([STORE_NAME], 'readonly')
      const store = transaction.objectStore(STORE_NAME)
      const request = store.get(id)
      
      request.onsuccess = function() {
        if (request.result) {
          resolve(request.result)
        } else {
          reject(new Error('模型不存在'))
        }
      }
      
      request.onerror = function() {
        reject(new Error('读取模型失败'))
      }
    })
  })
}

/**
 * 从 IndexedDB 删除模型数据
 * @param {string} id - 模型ID
 * @returns {Promise}
 */
export function deleteModelFromDB(id) {
  return openDB().then(db => {
    return new Promise((resolve, reject) => {
      const transaction = db.transaction([STORE_NAME], 'readwrite')
      const store = transaction.objectStore(STORE_NAME)
      const request = store.delete(id)
      
      request.onsuccess = function() {
        resolve()
      }
      
      request.onerror = function() {
        reject(new Error('删除模型失败'))
      }
    })
  })
}

/**
 * 清空所有模型数据
 * @returns {Promise}
 */
export function clearAllModelsFromDB() {
  return openDB().then(db => {
    return new Promise((resolve, reject) => {
      const transaction = db.transaction([STORE_NAME], 'readwrite')
      const store = transaction.objectStore(STORE_NAME)
      const request = store.clear()
      
      request.onsuccess = function() {
        resolve()
      }
      
      request.onerror = function() {
        reject(new Error('清空模型失败'))
      }
    })
  })
}
