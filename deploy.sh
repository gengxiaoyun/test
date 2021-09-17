#! /bin/bash

WORK_DIR=$(cd $(dirname $0); pwd)

INTERNET_AVAILABLE=false

PCRE_ARCHIVE_AVAILABLE=true
NGINX_ARCHIVE_AVAILABLE=true
GO_ARCHIVE_AVAILABLE=true
SOAR_ARCHIVE_AVAILABLE=true

PCRE_VERSION=8.35
PCRE_URL=http://downloads.sourceforge.net/project/pcre/pcre/${PCRE_VERSION}/pcre-${PCRE_VERSION}.tar.gz
PCRE_PATH=/data/pcre

NGINX_VERSION=1.21.1
NGINX_URL=http://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz
NGINX_PATH=/data/nginx

GO_VERSION=1.16.7
GO_URL=https://golang.google.cn/dl/go${GO_VERSION}.linux-amd64.tar.gz
OLD_GOROOT=$GOROOT
OLD_GOPATH=$GOPATH

SOAR_URL=https://github.com/romberli/soar.git
DAS_PATH=/data/das

function download() {
    for i in {1..10}
    do
        if [ ! -f $2 ]; then
            wget $1
        else
            return
        fi
    done
    echo "[ERROR] download $2 failed"
    exit 0
}


function deployDAS() {
    mkdir -p ${WORK_DIR}/archive
    
    mkdir -p ${DAS_PATH}/bin
    mkdir -p ${DAS_PATH}/conf
    
    checkInternet
    installDeps
    installNginx
    installGolang
    installDAS
    installSoar
}

function checkInternet() {
    # timeout
    local timeout=1
    # target url
    local target=www.baidu.com
    local retCode=`curl -I -s --connect-timeout ${timeout} ${target} -w %{http_code} | tail -n1`
    
    if [ "x$retCode" = "x200" ]; then
        INTERNET_AVAILABLE=true
        echo "[INFO] Network is available"
    else
        echo "[ERROR] Network is unavailable"
        exit 0
    fi
}

function installDeps() {
    yum -y install make zlib zlib-devel gcc-c++ libtool  openssl openssl-devel git #deps
}

function installNginx() {
    cd ${WORK_DIR}/archive
    
    download ${PCRE_URL} pcre-${PCRE_VERSION}.tar.gz
    download ${NGINX_URL} nginx-${NGINX_VERSION}.tar.gz
    
    tar -zxf pcre-${PCRE_VERSION}.tar.gz -C ${WORK_DIR}/archive
    tar -zxf nginx-${NGINX_VERSION}.tar.gz -C ${WORK_DIR}/archive
    
    mkdir -p /data
    
    mv ${WORK_DIR}/archive/pcre-${PCRE_VERSION} ${PCRE_PATH}
    cd ${PCRE_PATH}
    ./configure
    make && make install
    rm -rf ${WORK_DIR}/archive/pcre-${PCRE_VERSION}
    if [ ! -f ${PCRE_PATH}/pcre-config ]; then
        echo "[ERROR] install pcre failed"
        exit 0
    else
        echo "[INFO] install pcre success"
    fi
    
    cd ${WORK_DIR}/archive/nginx-${NGINX_VERSION}
    ./configure --prefix=${NGINX_PATH} \
    --conf-path=${NGINX_PATH}/conf/nginx.conf \
    --with-http_stub_status_module \
    --with-http_ssl_module \
    --with-pcre=${PCRE_PATH}
    make && make install
    rm -rf ${WORK_DIR}/archive/nginx-${NGINX_VERSION}
    if [ ! -f ${NGINX_PATH}/sbin/nginx ]; then
        echo "[ERROR] install nginx failed"
        exit 0
    else
        echo "[INFO] install nginx success"
    fi
    
    # mv config
    \cp -f ${WORK_DIR}/archive/nginx.conf ${NGINX_PATH}/conf
    
    echo "export NGINX_HOME=${NGINX_PATH}" >> /etc/profile
    echo 'export PATH=$PATH:$NGINX_HOME/sbin'  >> /etc/profile
}

function installGolang() {
    # judge if golang exist & version > 1.16
    local needGolang=true
    if [ -n "$GOROOT" ]; then # judge if $GOROOT exists
        local version=$(go version | grep -E 'go[1-9]\.((1[6-9])|([2-9][0-9]))')
        if [[ -n "$version" ]]; then
            needGolang=false
            echo "current golang meet the requirement of das"
        else
            echo "golang version less than 1.16.0"
        fi
        
    else
        echo "golang not found"
    fi
    
    # install golang
    if [ $needGolang = "true" ]; then
        cd ${WORK_DIR}/archive
        download ${GO_URL} go${GO_VERSION}.linux-amd64.tar.gz
        
        tar -zxf go${GO_VERSION}.linux-amd64.tar.gz -C ${WORK_DIR}/archive
        
        mv ${WORK_DIR}/archive/go /data
        rm -rf ${WORK_DIR}/archive/go
        
        # only current user will be affect with this setting
        export GOROOT=/data/go
        export GOPATH=${HOME}/go
        
        mkdir -p ${HOME}/go
        export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
    fi
}

function installSoar() {
    cd ${WORK_DIR}/../
    if [ ! -d ${WORK_DIR}/../soar ]; then
        git clone ${SOAR_URL}
    fi
    
    if [ -d ${WORK_DIR}/../soar ]; then
        cd soar
        make
    else
        echo "download soar failed"
        exit 0
    fi
    
    # mv das into /data
    # mv ${WORK_DIR}/../soar/bin/soar ${DAS_PATH}/bin
}

function installDAS() {
    cd ${WORK_DIR}
    echo "[INFO] compiling das.."
    make
    
    if [ $? -ne 0  ]; then
        echo "[ERROR] compiling das failed"
        exit 0
    else
        echo "[INFO] compilation success"
    fi
    
    # mv das into /data
    \cp -f ${WORK_DIR}/bin/das ${DAS_PATH}/bin
    \cp -f ${WORK_DIR}/config/das_default.yaml ${DAS_PATH}/conf/das.ymal
    
    # register to systemd
    chmod 0644 ${WORK_DIR}/archive/das.service
    \cp -f ${WORK_DIR}/archive/das.service /etc/systemd/system/
    systemctl enable das
    
    echo "export DAS_HOME=${DAS_PATH}"  >> /etc/profile
    echo 'export PATH=$PATH:$DAS_HOME/bin'  >> /etc/profile
}

deployDAS