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
          <el-radio :label="2">强制安装</el-radio>
          <el-radio :label="3">升级安装</el-radio>
          <el-radio :label="4">重新配置</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-divider content-position="left">V2ray配置选择</el-divider>
      <el-form-item label="选择配置">
        <el-radio-group v-model="form.config_type">
          <el-radio :label="1">预设配置(WebSocket/HTTPS)</el-radio>
          <el-radio :label="2">预设配置(KCP/HTTPS)</el-radio>
          <el-radio :label="3">自定义配置</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="监听端口" prop="port">
        <el-input v-model="form.v2ray_port" type="number" min="1" max="65535"
                  placeholder="V2ray服务器监听的端口号，默认值3000"></el-input>
      </el-form-item>
      <el-form-item label="传输方式" prop="transport_type">
        <el-select v-model="form.transport_type" class="w-100">
          <el-option :value="1" label="TCP"></el-option>
          <el-option :value="2" label="WebSocket"></el-option>
          <el-option :value="3" label="KCP"></el-option>
          <el-option :value="4" label="HTTP2"></el-option>
        </el-select>
      </el-form-item>
      <template v-for="client in form.clients" v-bind:key="client.user_id">
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
        <el-button type="primary" @click="addClient">添加新用户</el-button>
      </el-form-item>
      <template v-if="1 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">TCP传输配置</el-divider>
        <el-form-item label="伪装类型">
          <el-select v-model="form.tcp.type">
            <el-option value="none" label="NONE"></el-option>
            <el-option value="http" label="HTTP"></el-option>
          </el-select>
        </el-form-item>
        <template v-if="'http' === form.tcp.type">
          <div class="inline-form-item-2">
            <el-form-item label="请求版本" class="form-item">
              <el-select v-model="form.tcp.request.version">
                <el-option value="1.1" label="HTTP/1.1"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="请求方式" class="form-item">
              <el-select v-model="form.tcp.request.method">
                <el-option value="GET" label="GET"></el-option>
                <el-option value="POST" label="POST"></el-option>
              </el-select>
            </el-form-item>
          </div>
          <el-form-item label="请求路径">
            <el-input v-model="form.tcp.request.path" placeholder="请求路径，多个路径用英文,分隔，每次请求会随机选择一个，默认为/"></el-input>
          </el-form-item>
          <template v-for="header in form.tcp.request.headers" v-bind:key="header.key">
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
              <el-select v-model="form.tcp.response.version">
                <el-option value="1.1" label="HTTP/1.1"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="状态码" class="form-item">
              <el-input v-model="form.tcp.response.status" placeholder="HTTP响应的状态码，默认为200"></el-input>
            </el-form-item>
            <el-form-item label="状态描述" class="form-item">
              <el-input v-model="form.tcp.response.reason" placeholder="HTTP响应的状态描述，默认为OK"></el-input>
            </el-form-item>
          </div>
          <template v-for="header in form.tcp.response.headers" v-bind:key="header.key">
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
      <template v-if="2 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">WebSocket传输配置</el-divider>
        <el-form-item label="路径">
          <el-input v-model="form.web_socket.path" placeholder="URI路径"></el-input>
        </el-form-item>
        <template v-for="header in form.web_socket.headers" v-bind:key="header.key">
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
      <template v-if="3 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">KCP传输配置</el-divider>
        <div class="inline-form-item-2">
          <el-form-item label="伪装类型" class="form-item">
            <el-select v-model="form.kcp.type" class="w-100">
              <el-option value="none" label="none"></el-option>
              <el-option value="srtp" label="srtp"></el-option>
              <el-option value="utp" label="utp"></el-option>
              <el-option value="wechat-video" label="wechat-video"></el-option>
              <el-option value="dtls" label="dtls"></el-option>
              <el-option value="wireguard" label="wireguard"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="开启拥塞控制" class="form-item">
            <el-switch v-model="form.kcp.congestion"></el-switch>
          </el-form-item>
        </div>
        <div class="inline-form-item-2">
          <el-form-item label="MTU大小" class="form-item">
            <el-input v-model="form.kcp.mtu" type="number" min="576" max="1460"
                      placeholder="传输单元大小，576-1460之间的整数，默认为1350"></el-input>
          </el-form-item>
          <el-form-item label="TTI间隔时间" class="form-item">
            <el-input v-model="form.kcp.tti" type="number" min="10" max="100"
                      placeholder="传输间隔时间，10-100之间的整数，默认为50"></el-input>
          </el-form-item>
        </div>
        <div class="inline-form-item-2">
          <el-form-item label="上行带宽" class="form-item">
            <el-input v-model="form.kcp.uplink_capacity" type="number" min="0" placeholder="上行带宽大小，默认为5，单位MB/s">
              <template #append>MB/S</template>
            </el-input>
          </el-form-item>
          <el-form-item label="读取缓冲区大小" class="form-item">
            <el-input v-model="form.kcp.read_buffer_size" type="number" min="1" placeholder="读取缓冲区大小，默认为2，单位MB">
              <template #append>MB</template>
            </el-input>
          </el-form-item>
        </div>
        <div class="inline-form-item-2">
          <el-form-item label="下行带宽" class="form-item">
            <el-input v-model="form.kcp.downlink_capacity" type="number" min="0"
                      placeholder="下行带宽大小，默认为20，单位MB/s">
              <template #append>MB/S</template>
            </el-input>
          </el-form-item>
          <el-form-item label="写入缓冲区大小" class="form-item">
            <el-input v-model="form.kcp.write_buffer_size" type="number" min="1" placeholder="写入缓冲区大小，默认为2，单位MB">
              <template #append>MB</template>
            </el-input>
          </el-form-item>
        </div>
      </template>
      <template v-if="4 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">HTTP2传输配置</el-divider>
        <el-form-item label="域名">
          <el-input v-model="form.http2.host" placeholder="HTTP2的域名，多个使用英文,分隔"></el-input>
        </el-form-item>
        <el-form-item label="路径">
          <el-input v-model="form.http2.path" placeholder="URI路径"></el-input>
        </el-form-item>
      </template>
      <el-divider content-position="left">其他配置</el-divider>
      <el-form-item label="使用HTTPS">
        <el-switch v-model="form.use_tls"></el-switch>
      </el-form-item>
      <el-form-item label="HTTPS域名" prop="host" v-if="form.use_tls">
        <el-input v-model="form.tls_host" placeholder="HTTPS域名，该域名必须已被解析到目标服务器的IP地址"></el-input>
      </el-form-item>
      <el-form-item label-width="0" class="content-center">
        <el-button type="danger" @click="close">取消</el-button>
        <el-button type="primary" @click="save" :loading="saving">开始部署</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import {API_V2RAY_SERVER_DEPLOY, Client, Header, V2rayServerDeployForm} from "@/api/v2ray_server_develop"
import axios, {AxiosResponse} from "axios";
import {API_V2RAY_ENDPOINT_NEW} from "@/api/v2ray_endpoint_new";
import {BaseResponse} from "@/api/base";

export default defineComponent({
  name: "DeployV2rayServer",

  props: {
    show: Boolean,
  },

  emits: ['update:show', 'success'],

  watch: {
    show: function () {
      this.form.v2ray_port = 3000
      this.form.transport_type = 2
      this.form.web_socket.path = "/example-path"
      this.form.use_tls = true

      this.form.clients = new Array<Client>()
      this.addClient()
    },

    'form.config_type': function () {
      if (1 == this.form.config_type) {
        this.form.v2ray_port = 3000
        this.form.transport_type = 2
        this.form.web_socket.path = "/example-path"
        this.form.use_tls = true
      }

      if (2 == this.form.config_type) {
        this.form.v2ray_port = 3000
        this.form.transport_type = 3
        this.form.use_tls = true
      }

      if (3 == this.form.config_type) {
        this.form.v2ray_port = 3000
        this.form.transport_type = 1
        this.form.use_tls = false
      }
    }
  },

  data() {
    return {
      saving: false,
      form: new V2rayServerDeployForm(),
      rules: {
        // host: [{required: true, message: '必须填写主机地址', trigger: 'blur'}],
        // port: [{required: true, message: '必须填写端口号', trigger: 'blur'}],
        // user_id: [{required: true, message: '必须填写用户ID', trigger: 'blur'}],
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

        this.form.v2ray_port = parseInt(this.form.v2ray_port.toString())

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
      this.form.web_socket.headers.push(new Header())
    },

    addTcpRequestHeader() {
      this.form.tcp.request.headers.push(new Header())
    },

    addTcpResponseHeader() {
      this.form.tcp.response.headers.push(new Header())
    },

    addClient() {
      this.form.clients.push(new Client())
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
