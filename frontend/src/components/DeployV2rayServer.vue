<template>
  <el-dialog :close-on-click-modal="false" :close-on-press-escape="false" :model-value="show" destroy-on-close
             width="80%" @close="close">
    <el-form ref="V2rayServerDeploy" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="安装方式">
        <el-radio-group v-model="form.install_type">
          <el-radio :label="1">默认安装</el-radio>
          <el-radio :label="2">重新安装</el-radio>
          <el-radio :label="3">仅升级V2ray</el-radio>
          <el-radio :label="4">仅配置V2ray</el-radio>
        </el-radio-group>
      </el-form-item>
      <template v-if="3 !== parseInt(form.install_type.toString())">
        <el-divider content-position="left">V2ray配置选择</el-divider>
        <el-form-item label="选择配置">
          <el-radio-group v-model="form.config_type">
            <el-radio :label="1">预设配置(WebSocket)</el-radio>
            <el-radio :label="2">预设配置(KCP)</el-radio>
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
        <el-form-item
            v-if="1 === parseInt(form.v2ray_config.transport_type.toString()) || 3 === parseInt(form.v2ray_config.transport_type.toString())"
            label="监听端口">
          <el-input v-model="form.v2ray_config.v2ray_port" placeholder="V2ray监听的端口号"></el-input>
        </el-form-item>
        <template v-for="client in form.v2ray_config.clients" v-bind:key="client">
          <el-form-item label="用户ID">
            <el-input v-model="client.user_id" placeholder="用户ID，请勿使用过短的用户ID，若不填写将会自动生成"></el-input>
          </el-form-item>
        </template>
        <el-form-item class="content-center" label-width="0">
          <el-button type="primary" @click="addClient">添加一个用户</el-button>
        </el-form-item>
        <template v-if="1 === parseInt(form.v2ray_config.transport_type.toString())">
          <el-divider content-position="left">TCP传输配置</el-divider>
          <el-form-item label="伪装类型">
            <el-select v-model="form.v2ray_config.tcp.type">
              <el-option label="NONE" value="none"></el-option>
              <el-option label="HTTP" value="http"></el-option>
            </el-select>
          </el-form-item>
          <template v-if="'http' === form.v2ray_config.tcp.type">
            <div class="inline-form-item-2">
              <el-form-item class="form-item" label="请求版本">
                <el-select v-model="form.v2ray_config.tcp.request.version">
                  <el-option label="HTTP/1.1" value="1.1"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item class="form-item" label="请求方式">
                <el-select v-model="form.v2ray_config.tcp.request.method">
                  <el-option label="GET" value="GET"></el-option>
                  <el-option label="POST" value="POST"></el-option>
                </el-select>
              </el-form-item>
            </div>
            <el-form-item label="请求路径">
              <el-input v-model="form.v2ray_config.tcp.request.path"
                        placeholder="请求路径，多个路径用英文,分隔，每次请求会随机选择一个，默认为/"></el-input>
            </el-form-item>
            <template v-for="header in form.v2ray_config.tcp.request.headers" v-bind:key="header">
              <div class="inline-form-item-2">
                <el-form-item class="form-item" label="字段名">
                  <el-input v-model="header.key" placeholder="自定义头的字段名称"></el-input>
                </el-form-item>
                <el-form-item class="form-item" label="字段值">
                  <el-input v-model="header.value" placeholder="自定义头的字段值"></el-input>
                </el-form-item>
              </div>
            </template>
            <el-form-item class="content-center" label-width="0">
              <el-button type="primary" @click="addTcpRequestHeader">添加自定义请求头字段</el-button>
            </el-form-item>
            <div class="inline-form-item-3">
              <el-form-item class="form-item" label="响应版本">
                <el-select v-model="form.v2ray_config.tcp.response.version">
                  <el-option label="HTTP/1.1" value="1.1"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item class="form-item" label="状态码">
                <el-input v-model="form.v2ray_config.tcp.response.status" placeholder="HTTP响应的状态码，默认为200"></el-input>
              </el-form-item>
              <el-form-item class="form-item" label="状态描述">
                <el-input v-model="form.v2ray_config.tcp.response.reason" placeholder="HTTP响应的状态描述，默认为OK"></el-input>
              </el-form-item>
            </div>
            <template v-for="header in form.v2ray_config.tcp.response.headers" v-bind:key="header">
              <div class="inline-form-item-2">
                <el-form-item class="form-item" label="字段名">
                  <el-input v-model="header.key" placeholder="自定义头的字段名称"></el-input>
                </el-form-item>
                <el-form-item class="form-item" label="字段值">
                  <el-input v-model="header.value" placeholder="自定义头的字段值"></el-input>
                </el-form-item>
              </div>
            </template>
            <el-form-item class="content-center" label-width="0">
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
              <el-form-item class="form-item" label="字段名">
                <el-input v-model="header.key" placeholder="自定义头的字段名称"></el-input>
              </el-form-item>
              <el-form-item class="form-item" label="字段值">
                <el-input v-model="header.value" placeholder="自定义头的字段值"></el-input>
              </el-form-item>
            </div>
          </template>
          <el-form-item class="content-center" label-width="0">
            <el-button type="primary" @click="addWebSocketHeader">添加自定义头字段</el-button>
          </el-form-item>
        </template>
        <template v-if="3 === parseInt(form.v2ray_config.transport_type.toString())">
          <el-divider content-position="left">KCP传输配置</el-divider>
          <div class="inline-form-item-2">
            <el-form-item class="form-item" label="伪装类型">
              <el-select v-model="form.v2ray_config.kcp.type" class="w-100">
                <el-option label="none" value="none"></el-option>
                <el-option label="srtp" value="srtp"></el-option>
                <el-option label="utp" value="utp"></el-option>
                <el-option label="wechat-video" value="wechat-video"></el-option>
                <el-option label="dtls" value="dtls"></el-option>
                <el-option label="wireguard" value="wireguard"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item class="form-item" label="开启拥塞控制">
              <el-switch v-model="form.v2ray_config.kcp.congestion"></el-switch>
            </el-form-item>
          </div>
          <div class="inline-form-item-2">
            <el-form-item class="form-item" label="MTU大小">
              <el-input v-model="form.v2ray_config.kcp.mtu" max="1460" min="576"
                        placeholder="传输单元大小，576-1460之间的整数，默认为1350"
                        type="number"></el-input>
            </el-form-item>
            <el-form-item class="form-item" label="TTI间隔时间">
              <el-input v-model="form.v2ray_config.kcp.tti" max="100" min="10" placeholder="传输间隔时间，10-100之间的整数，默认为50"
                        type="number"></el-input>
            </el-form-item>
          </div>
          <div class="inline-form-item-2">
            <el-form-item class="form-item" label="上行带宽">
              <el-input v-model="form.v2ray_config.kcp.uplink_capacity" min="0" placeholder="上行带宽大小，默认为5，单位MB/s"
                        type="number">
                <template #append>MB/S</template>
              </el-input>
            </el-form-item>
            <el-form-item class="form-item" label="读取缓冲区大小">
              <el-input v-model="form.v2ray_config.kcp.read_buffer_size" min="1" placeholder="读取缓冲区大小，默认为2，单位MB"
                        type="number">
                <template #append>MB</template>
              </el-input>
            </el-form-item>
          </div>
          <div class="inline-form-item-2">
            <el-form-item class="form-item" label="下行带宽">
              <el-input v-model="form.v2ray_config.kcp.downlink_capacity" min="0" placeholder="下行带宽大小，默认为20，单位MB/s"
                        type="number">
                <template #append>MB/S</template>
              </el-input>
            </el-form-item>
            <el-form-item class="form-item" label="写入缓冲区大小">
              <el-input v-model="form.v2ray_config.kcp.write_buffer_size" min="1" placeholder="写入缓冲区大小，默认为2，单位MB"
                        type="number">
                <template #append>MB</template>
              </el-input>
            </el-form-item>
          </div>
        </template>
        <template v-if="4 === parseInt(form.v2ray_config.transport_type.toString())">
          <el-divider content-position="left">HTTP2传输配置</el-divider>
          <el-form-item label="域名">
            <el-input v-model="form.v2ray_config.http2.host"
                      placeholder="HTTP2的域名列表，多个使用英文,分隔。列表内会自动添加HTTPS的域名，请不要重复添加"></el-input>
          </el-form-item>
          <el-form-item label="路径">
            <el-input v-model="form.v2ray_config.http2.path" placeholder="URI路径"></el-input>
          </el-form-item>
        </template>
        <template v-if="4 !== parseInt(form.install_type.toString())">
          <el-divider content-position="left">HTTPS配置
            <el-tooltip content="除了KCP模式都支持HTTPS协议，V2rayHelper可以进行自动证书申请及续期" placement="right"><i
                class="el-icon-info"></i></el-tooltip>
          </el-divider>
          <el-form-item label="使用HTTPS">
            <el-switch v-model="form.use_tls"
                       :disabled="3 === parseInt(form.v2ray_config.transport_type.toString()) || 4 === parseInt(form.v2ray_config.transport_type.toString())"></el-switch>
          </el-form-item>
          <el-form-item v-if="form.use_tls" label="HTTPS域名" prop="tls_host">
            <el-input v-model="form.tls_host" placeholder="HTTPS域名，该域名必须已被解析到目标服务器的IP地址"></el-input>
          </el-form-item>
          <el-divider content-position="left">站点伪装配置
            <el-tooltip content="可以在部署V2ray的同时部署一个站点，可以增强伪装效果。TCP模式与KCP模式不支持站点伪装" placement="right">
              <i class="el-icon-info"></i>
            </el-tooltip>
          </el-divider>
          <div class="inline-form-item-2">
            <el-form-item class="form-item-0" label="开启伪装">
              <el-switch v-model="form.enable_web_service"
                         :disabled="1 === parseInt(form.v2ray_config.transport_type.toString()) || 3 === parseInt(form.v2ray_config.transport_type.toString())"></el-switch>
            </el-form-item>
            <el-form-item class="form-item-1" label="伪装站点" label-width="80px">
              <el-select v-model="form.web_service_type" class="w-100">
                <el-option value="cloudreve" label="Cloudreve"></el-option>
              </el-select>
            </el-form-item>
          </div>
          <template v-if="form.enable_web_service && 'cloudreve' === form.web_service_type">
            <el-divider content-position="left">Cloudreve配置信息</el-divider>
            <el-form-item label="Cloudreve配置">
              <el-checkbox v-model="form.cloudreve_config.enable_aria2" label="启用Aria2离线下载支持(不懂不要选)"></el-checkbox>
              <el-checkbox v-model="form.cloudreve_config.reset_admin_password"
                           label="重置管理员密码(首次部署请选中，否则无法获取初始密码)"></el-checkbox>
            </el-form-item>
            <div class="inline-form-item-2" v-if="form.cloudreve_config.reset_admin_password">
              <el-form-item class="form-item" label="初始管理员账号">
                <el-input v-model="response.cloudreve_admin" readonly placeholder="将在部署成功后回显"></el-input>
              </el-form-item>
              <el-form-item class="form-item" label="初始管理员密码">
                <el-input v-model="response.cloudreve_password" readonly placeholder="将在部署成功后回显"></el-input>
              </el-form-item>
            </div>
          </template>
        </template>
      </template>
      <el-form-item class="content-center" label-width="0">
        <el-button type="danger" @click="close" :disabled="deploying">关闭</el-button>
        <el-button :loading="deploying" type="primary" @click="save">开始部署</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import {
  API_V2RAY_SERVER_DEPLOY,
  Client,
  V2rayServerDeployData,
  V2rayServerDeployForm,
  V2rayServerDeployResponse
} from "@/api/v2ray_server_develop"
import axios, {AxiosResponse} from "axios"
import {Header} from "@/api/base"

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
      this.form.enable_web_service = true

      this.form.v2ray_config.clients = new Array<Client>()
      this.addClient()
    },

    'form.config_type': function () {
      if (1 == this.form.config_type) {
        this.form.v2ray_config.v2ray_port = 3000
        this.form.v2ray_config.transport_type = 2
        this.form.v2ray_config.web_socket.path = "/v2ray-ws-path"
        this.form.use_tls = true
        this.form.enable_web_service = true
      }

      if (2 == this.form.config_type) {
        this.form.v2ray_config.v2ray_port = 3000
        this.form.v2ray_config.transport_type = 3
        this.form.use_tls = false
        this.form.enable_web_service = false
      }

      if (3 == this.form.config_type) {
        this.form.v2ray_config.v2ray_port = 3000
        this.form.v2ray_config.transport_type = 1
        this.form.use_tls = false
        this.form.enable_web_service = false
      }
    },

    'form.v2ray_config.transport_type': function () {
      if (4 == this.form.v2ray_config.transport_type) {
        this.form.use_tls = true
      }

      if (3 == this.form.v2ray_config.transport_type) {
        this.form.use_tls = false
      }

      if (1 == this.form.v2ray_config.transport_type || 3 == this.form.v2ray_config.transport_type) {
        this.form.enable_web_service = false
      }
    }
  },

  data() {
    return {
      deploying: false,
      form: new V2rayServerDeployForm(),
      rules: {
        tls_host: [{validator: this.validateTlsHost, trigger: 'blur'}],
      },
      headers: [],
      response: new V2rayServerDeployData()
    }
  },

  methods: {
    close() {
      this.form = new V2rayServerDeployForm()

      this.$emit('update:show', false)
    },

    success() {
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

        this.deploying = true
        axios.post(API_V2RAY_SERVER_DEPLOY, this.form).then((response: AxiosResponse<V2rayServerDeployResponse>) => {
          this.deploying = false
          if (0 != response.data.code) {
            this.$message.error(response.data.message)

            return
          }

          this.response = response.data.data
          this.$message.success("部署成功，已自动生成配置文件")
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
  },
})
</script>

<style lang="scss" scoped>
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
