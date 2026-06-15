import { request, METHOD } from '@/utils/request'

import {
    AI2DGENERATEIMAGE
} from '@/services/api'

export async function GenerateISM2DScene(params) {
  return request(AI2DGENERATEIMAGE, METHOD.POST, params, {
    timeout: 120000
  })
}

export default {
  GenerateISM2DScene
}
