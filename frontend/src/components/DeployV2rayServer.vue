<template>
  <el-dialog width="60%" :model-value="show" @close="close" :close-on-click-modal="false"
             :close-on-press-escape="false" destroy-on-close>
    <el-form :model="form" label-width="120px" ref="V2rayServerDeploy" :rules="rules">
      <el-form-item label="选择服务器">
        <el-radio-group v-model="form.server_type">
          <el-radio :label="1">当前服务器</el-radio>
          <el-radio :label="2" disabled>远程服务器(暂未支持)</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="安装方式">
        <el-radio-group v-model="form.install_type">
          <el-radio :label="1">默认安装</el-radio>
          <el-radio :label="2">重新安装</el-radio>
          <el-radio :label="3">仅升级V2ray</el-radio>
          <el-radio :label="4">仅配置V2ray</el-radio>
        </el-radio-group>
      </el-form-item>
      <template v-if="3 !== parseInt(form.install_type.toString()) && 5 !== parseInt(form.install_type.toString())">
        <el-divider content-position="left">V2ray配置选择</el-divider>
        <el-form-item label="选择配置">
          <el-radio-group v-model="form.config_type">
            <el-radio :label="1">预设配置(WebSocket/HTTPS)</el-radio>
            <el-radio :label="2">预设配置(KCP/HTTPS)</el-radio>
            <el-radio :label="3">自定义配置</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="传输方式" prop="transport_type">
          <el-select v-model="form.v2ray_config.transport_type" class="w-100">
            <el-option :value="1" label="TCP"></el-option>
            <el-option :value="2" label="WebSocket"></el-option>
            <el-option :value="3" label="KCP"></el-option>
            <el-option :value="4" label="HTTP2"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="监听端口"
                      v-if="1 === parseInt(form.v2ray_config.transport_type.toString()) || 3 === parseInt(form.v2ray_config.transport_type.toString())">
          <el-input v-model="form.v2ray_config.v2ray_port" placeholder="V2ray监听的端口号"></el-input>
        </el-form-item>
        <template v-for="client in form.v2ray_config.clients" v-bind:key="client.user_id">
          <div class="inline-form-item-client">
            <el-form-item label="用户ID" class="form-item-user-id">
              <el-input v-model="client.user_id" placeholder="用户ID，请勿使用过短的用户ID，若不填写将会自动生成"></el-input>
            </el-form-item>
            <el-form-item label="额外ID" label-width="60px" class="form-item-alter-id">
              <el-input v-model="client.alter_id" type="number" min="0" max="65535"
                        placeholder="额外ID的数量，推荐值4，可选值0-65535"></el-input>
            </el-form-item>
          </div>
        </template>
        <el-form-item label-width="0" class="content-center">
          <el-button type="primary" @click="addClient">添加一个用户</el-button>
        </el-form-item>
        <template v-if="1 === parseInt(form.v2ray_config.transport_type.toString())">
          <el-divider content-position="left">TCP传输配置</el-divider>
          <el-form-item label="伪装类型">
            <el-select v-model="form.v2ray_config.tcp.type">
              <el-option value="none" label="NONE"></el-option>
              <el-option value="http" label="HTTP"></el-option>
            </el-select>
          </el-form-item>
          <template v-if="'http' === form.v2ray_config.tcp.type">
            <div class="inline-form-item-2">
              <el-form-item label="请求版本" class="form-item">
                <el-select v-model="form.v2ray_config.tcp.request.version">
                  <el-option value="1.1" label="HTTP/1.1"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="请求方式" class="form-item">
                <el-select v-model="form.v2ray_config.tcp.request.method">
                  <el-option value="GET" label="GET"></el-option>
                  <el-option value="POST" label="POST"></el-option>
                </el-select>
              </el-form-item>
            </div>
            <el-form-item label="请求路径">
              <el-input v-model="form.v2ray_config.tcp.request.path"
                        placeholder="请求路径，多个路径用英文,分隔，每次请求会随机选择一个，默认为/"></el-input>
            </el-form-item>
            <template v-for="header in form.v2ray_config.tcp.request.headers" v-bind:key="header">
              <div class="inline-form-item-2">
                <el-form-item label="字段名" class="form-item">
                  <el-input v-model="header.key" placeholder="自定义头的字段名称"></el-input>
                </el-form-item>
                <el-form-item label="字段值" class="form-item">
                  <el-input v-model="header.value" placeholder="自定义头的字段值"></el-input>
                </el-form-item>
              </div>
            </template>
            <el-form-item label-width="0" class="content-center">
              <el-button type="primary" @click="addTcpRequestHeader">添加自定义请求头字段</el-button>
            </el-form-item>
            <div class="inline-form-item-3">
              <el-form-item label="响应版本" class="form-item">
                <el-select v-model="form.v2ray_config.tcp.response.version">
                  <el-option value="1.1" label="HTTP/1.1"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="状态码" class="form-item">
                <el-input v-model="form.v2ray_config.tcp.response.status" placeholder="HTTP响应的状态码，默认为200"></el-input>
              </el-form-item>
              <el-form-item label="状态描述" class="form-item">
                <el-input v-model="form.v2ray_config.tcp.response.reason" placeholder="HTTP响应的状态描述，默认为OK"></el-input>
              </el-form-item>
            </div>
            <template v-for="header in form.v2ray_config.tcp.response.headers" v-bind:key="header">
              <div class="inline-form-item-2">
                <el-form-item label="字段名" class="form-item">
                  <el-input v-model="header.key" placeholder="自定义头的字段名称"></el-input>
                </el-form-item>
                <el-form-item label="字段值" class="form-item">
                  <el-input v-model="header.value" placeholder="自定义头的字段值"></el-input>
                </el-form-item>
              </div>
            </template>
            <el-form-item label-width="0" class="content-center">
              <el-button type="primary" @click="addTcpResponseHeader">添加自定义响应头字段</el-button>
            </el-form-item>
          </template>
        </template>
        <template v-if="2 === parseInt(form.v2ray_config.transport_type.toString())">
          <el-divider content-position="left">WebSocket传输配置</el-divider>
          <el-form-item label="路径">
            <el-input v-model="form.v2ray_config.web_socket.path" placeholder="URI路径"></el-input>
          </el-form-item>
          <template v-for="header in form.v2ray_config.web_socket.headers" v-bind:key="header">
            <div class="inline-form-item-2">
              <el-form-item label="字段名" class="form-item">
                <el-input v-model="header.key" placeholder="自定义头的字段名称"></el-input>
              </el-form-item>
              <el-form-item label="字段值" class="form-item">
                <el-input v-model="header.value" placeholder="自定义头的字段值"></el-input>
              </el-form-item>
            </div>
          </template>
          <el-form-item label-width="0" class="content-center">
            <el-button type="primary" @click="addWebSocketHeader">添加自定义头字段</el-button>
          </el-form-item>
        </template>
        <template v-if="3 === parseInt(form.v2ray_config.transport_type.toString())">
          <el-divider content-position="left">KCP传输配置</el-divider>
          <div class="inline-form-item-2">
            <el-form-item label="伪装类型" class="form-item">
              <el-select v-model="form.v2ray_config.kcp.type" class="w-100">
                <el-option value="none" label="none"></el-option>
                <el-option value="srtp" label="srtp"></el-option>
                <el-option value="utp" label="utp"></el-option>
                <el-option value="wechat-video" label="wechat-video"></el-option>
                <el-option value="dtls" label="dtls"></el-option>
                <el-option value="wireguard" label="wireguard"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="开启拥塞控制" class="form-item">
              <el-switch v-model="form.v2ray_config.kcp.congestion"></el-switch>
            </el-form-item>
          </div>
          <div class="inline-form-item-2">
            <el-form-item label="MTU大小" class="form-item">
              <el-input v-model="form.v2ray_config.kcp.mtu" type="number" min="576" max="1460"
                        placeholder="传输单元大小，576-1460之间的整数，默认为1350"></el-input>
            </el-form-item>
            <el-form-item label="TTI间隔时间" class="form-item">
              <el-input v-model="form.v2ray_config.kcp.tti" type="number" min="10" max="100"
                        placeholder="传输间隔时间，10-100之间的整数，默认为50"></el-input>
            </el-form-item>
          </div>
          <div class="inline-form-item-2">
            <el-form-item label="上行带宽" class="form-item">
              <el-input v-model="form.v2ray_config.kcp.uplink_capacity" type="number" min="0"
                        placeholder="上行带宽大小，默认为5，单位MB/s">
                <template #append>MB/S</template>
              </el-input>
            </el-form-item>
            <el-form-item label="读取缓冲区大小" class="form-item">
              <el-input v-model="form.v2ray_config.kcp.read_buffer_size" type="number" min="1"
                        placeholder="读取缓冲区大小，默认为2，单位MB">
                <template #append>MB</template>
              </el-input>
            </el-form-item>
          </div>
          <div class="inline-form-item-2">
            <el-form-item label="下行带宽" class="form-item">
              <el-input v-model="form.v2ray_config.kcp.downlink_capacity" type="number" min="0"
                        placeholder="下行带宽大小，默认为20，单位MB/s">
                <template #append>MB/S</template>
              </el-input>
            </el-form-item>
            <el-form-item label="写入缓冲区大小" class="form-item">
              <el-input v-model="form.v2ray_config.kcp.write_buffer_size" type="number" min="1"
                        placeholder="写入缓冲区大小，默认为2，单位MB">
                <template #append>MB</template>
              </el-input>
            </el-form-item>
          </div>
        </template>
        <template
            v-if="4 !== parseInt(form.install_type.toString()) && 3 !== parseInt(form.v2ray_config.transport_type.toString())">
          <el-divider content-position="left">HTTPS配置
            <el-tooltip content="WebSocket模式通过Caddy自动申请及续期证书；TCP模式可以自动申请证书或手动上传证书" placement="right"><i
                class="el-icon-info"></i></el-tooltip>
          </el-divider>
          <el-form-item label="使用HTTPS">
            <el-switch v-model="form.use_tls"></el-switch>
          </el-form-item>
          <el-form-item label="HTTPS域名" prop="tls_host" v-if="form.use_tls">
            <el-input v-model="form.tls_host" placeholder="HTTPS域名，该域名必须已被解析到目标服务器的IP地址"></el-input>
          </el-form-item>
          <el-form-item label="HTTPS证书" v-if="form.use_tls">
            <el-radio-group v-model="form.cert_type">
              <el-radio :label="1">自动申请证书</el-radio>
              <el-radio :label="2" disabled>上传证书(暂未支持)</el-radio>
            </el-radio-group>
          </el-form-item>
          <!--          <el-divider content-position="left">Cloudreve配置-->
          <!--            <el-tooltip content="Cloudreve是一个轻量级私有云网盘服务，同时布置可增强伪装效果" placement="right"><i class="el-icon-info"></i>-->
          <!--            </el-tooltip>-->
          <!--          </el-divider>-->
          <!--          <el-form-item label="配置Cloudreve" prop="use_cloudreve">-->
          <!--            <el-switch v-model="form.use_cloudreve"></el-switch>-->
          <!--          </el-form-item>-->
        </template>
      </template>
      <el-form-item label-width="0" class="content-center">
        <el-button type="danger" @click="close">取消</el-button>
        <el-button type="primary" @click="save" :loading="saving">开始部署</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import {API_V2RAY_SERVER_DEPLOY, Client, V2rayServerDeployForm} from "@/api/v2ray_server_develop"
import axios, {AxiosResponse} from "axios"
import {BaseResponse, Header} from "@/api/base"

export default defineComponent({
  name: "DeployV2rayServer",

  props: {
    show: Boolean,
  },

  emits: ['update:show', 'success'],

  watch: {
    show: function () {
      this.form.v2ray_config.v2ray_port = 3000
      this.form.v2ray_config.transport_type = 2
      this.form.v2ray_config.web_socket.path = "/v2ray-ws-path"
      this.form.use_tls = true
      this.form.use_cloudreve = false

      this.form.v2ray_config.clients = new Array<Client>()
      this.addClient()
    },

    'form.config_type': function () {
      if (1 == this.form.config_type) {
        this.form.v2ray_config.v2ray_port = 3000
        this.form.v2ray_config.transport_type = 2
        this.form.v2ray_config.web_socket.path = "/v2ray-ws-path"
        this.form.use_tls = true
        this.form.use_cloudreve = false
      }

      if (2 == this.form.config_type) {
        this.form.v2ray_config.v2ray_port = 3000
        this.form.v2ray_config.transport_type = 3
        this.form.use_tls = false
        this.form.use_cloudreve = false
      }

      if (3 == this.form.config_type) {
        this.form.v2ray_config.v2ray_port = 3000
        this.form.v2ray_config.transport_type = 1
        this.form.use_tls = false
        this.form.use_cloudreve = false
      }
    }
  },

  data() {
    return {
      saving: false,
      form: new V2rayServerDeployForm(),
      rules: {
        tls_host: [{validator: this.validateTlsHost, trigger: 'blur'}],
        use_cloudreve: [{validator: this.validateUseCloudreve, trigger: 'blur'}],
      },
      headers: [],
    }
  },

  methods: {
    close() {
      this.form = new V2rayServerDeployForm()

      this.$emit('update:show', false)
    },

    success() {
      this.close()

      this.$emit('success')
    },

    save() {
      let form = this.$refs.V2rayServerDeploy as any

      form.validate((isValid: boolean) => {
        if (!isValid) {
          return
        }

        this.form.v2ray_config.v2ray_port = parseInt(this.form.v2ray_config.v2ray_port.toString())

        this.form.v2ray_config.kcp.mtu = parseInt(this.form.v2ray_config.kcp.mtu.toString())
        this.form.v2ray_config.kcp.tti = parseInt(this.form.v2ray_config.kcp.tti.toString())
        this.form.v2ray_config.kcp.uplink_capacity = parseInt(this.form.v2ray_config.kcp.uplink_capacity.toString())
        this.form.v2ray_config.kcp.downlink_capacity = parseInt(this.form.v2ray_config.kcp.downlink_capacity.toString())
        this.form.v2ray_config.kcp.read_buffer_size = parseInt(this.form.v2ray_config.kcp.read_buffer_size.toString())
        this.form.v2ray_config.kcp.write_buffer_size = parseInt(this.form.v2ray_config.kcp.write_buffer_size.toString())

        for (let i = 0; i < this.form.v2ray_config.clients.length; i++) {
          this.form.v2ray_config.clients[i].alter_id = parseInt(this.form.v2ray_config.clients[i].alter_id.toString())
        }

        this.saving = true
        axios.post(API_V2RAY_SERVER_DEPLOY, this.form).then((response: AxiosResponse<BaseResponse>) => {
          this.saving = false
          if (0 != response.data.code) {
            this.$message.error(response.data.message)

            return
          }

          this.success()
        })
      })
    },

    addWebSocketHeader() {
      this.form.v2ray_config.web_socket.headers.push(new Header())
    },

    addTcpRequestHeader() {
      this.form.v2ray_config.tcp.request.headers.push(new Header())
    },

    addTcpResponseHeader() {
      this.form.v2ray_config.tcp.response.headers.push(new Header())
    },

    addClient() {
      this.form.v2ray_config.clients.push(new Client())
    },

    validateTlsHost(a: any, b: any, c: any) {
      if (this.form.use_tls && "" == this.form.tls_host.trim()) {
        c(new Error("开启HTTPS时必须填写HTTPS域名"))

        return
      }

      c()
    },

    validateUseCloudreve(a: any, b: any, c: any) {
      if (this.form.use_cloudreve && (1 == this.form.v2ray_config.transport_type || 3 == this.form.v2ray_config.transport_type)) {
        c(new Error("配置Cloudreve时V2ray不支持使用TCP与KCP，请使用WebSocket或HTTP2"))

        return
      }

      c()
    }
  },
})
</script>

<style scoped lang="scss">
.inline-form-item-2 {
  display: flex;

  .form-item {
    width: 50%;
  }
}

.inline-form-item-3 {
  display: flex;

  .form-item {
    width: 33.33%;
  }
}

.inline-form-item-client {
  display: flex;

  .form-item-user-id {
    width: 60%;
  }

  .form-item-alter-id {
    width: 40%;
  }
}
</style>
