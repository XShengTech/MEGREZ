import axios from "axios";

const ajax = (url, method, { params = {}, data = {} }) => {
  axios.defaults.withCredentials = true
  axios.defaults.crossDomain = true
  axios.defaults.baseURL = '/api/v1'
  return new Promise((resolve, reject) => {
    axios({
      url,
      method,
      params,
      data
    }).then(res => {
      if (res.data.code != 200) {
        reject(res)
      } else {
        resolve(res)
      }
    }, res => {
      reject(res)
    })
  })
}

export default {
  GetUserProfile() {
    return ajax('user/profile', 'get', {})
  },
  UserLogin(data) {
    return ajax('user/login', 'post', { data })
  },
  UserLogout() {
    return ajax('user/logout', 'get', {})
  },
  UserRegister(data) {
    return ajax('user/register', 'post', { data })
  },
  UserResetPassword(data) {
    return ajax('user/password', 'post', { data })
  },
  UserVerifyRequest() {
    return ajax(`user/verify`, 'post', {})
  },
  UserVerify(code) {
    return ajax(`user/verify/${code}`, 'get', {})
  },

  UserInstancesList(params) {
    return ajax('user/instances', 'get', { params })
  },
  UserInstancesDetail(id) {
    return ajax(`user/instances/${id}`, 'get', {})
  },
  UserInstancesModify(id, data) {
    return ajax(`user/instances/${id}`, 'post', { data })
  },
  UserInstancesModifyLabel(id, data) {
    return ajax(`user/instances/${id}/label`, 'post', { data })
  },
  UserInstancesAction(id, data) {
    return ajax(`user/instances/${id}/`, 'put', { data })
  },
  UserInstancesCreate(data) {
    return ajax('user/instances', 'post', { data })
  },
  UserInstancesDelete(id) {
    return ajax(`user/instances/${id}`, 'delete', {})
  },

  UserServerList(params) {
    return ajax('user/servers', 'get', { params })
  },
  UserServerDetail(id) {
    return ajax(`user/servers/${id}`, 'get', {})
  },

  UserImages() {
    return ajax('user/images', 'get', {})
  },

  AdminInstancesList(params) {
    return ajax('admin/instances', 'get', { params })
  },
  AdminInstancesDetail(id) {
    return ajax(`admin/instances/${id}`, 'get', {})
  },
  AdminInstancesModify(id, data) {
    return ajax(`admin/instances/${id}`, 'post', { data })
  },
  AdminInstancesModifyLabel(id, data) {
    return ajax(`admin/instances/${id}/label`, 'post', { data })
  },
  AdminInstancesAction(id, data) {
    return ajax(`admin/instances/${id}/`, 'put', { data })
  },
  AdminInstancesCreate(data) {
    return ajax('admin/instances', 'post', { data })
  },
  AdminInstancesDelete(id) {
    return ajax(`admin/instances/${id}`, 'delete', {})
  },
  AdminInstancesForceDelete(id) {
    return ajax(`admin/instances/${id}/force`, 'delete', {})
  },

  AdminUserList(params) {
    return ajax('admin/users', 'get', { params })
  },
  AdminUserModify(id, data) {
    return ajax(`admin/users/${id}`, 'post', { data })
  },
  AdminUserDelete(id) {
    return ajax(`admin/users/${id}`, 'delete', {})
  },

  AdminServersList(params) {
    return ajax('admin/servers', 'get', { params })
  },
  AdminServersDetail(id) {
    return ajax(`admin/servers/${id}`, 'get', {})
  },
  AdminServersAdd(data) {
    return ajax('admin/servers', 'post', { data })
  },
  AdminServersModify(id, data) {
    return ajax(`admin/servers/${id}`, 'post', { data })
  },
  AdminServersDelete(id) {
    return ajax(`admin/servers/${id}`, 'delete', {})
  },

  AdminImagesList() {
    return ajax('admin/images', 'get', {})
  },
  AdminImagesModify(data) {
    return ajax('admin/images', 'post', { data })
  },
}