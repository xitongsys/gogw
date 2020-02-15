MSG_FRESH_CHART = "fresh_chart_"

function SpeedToReadable(s){
    if(s <= 1024){
        return s + " B/s"
    }else if(s < 1024 * 1024){
        return (s/1024.0).toFixed(2) + " KB/s"
    }else{
        return (s/1024.0/1024.0).toFixed(2) + " MB/s"
    }
}

function Client(divid){
    return {
        DivId: divid,
        ClientId: "0",
        ClientAddr: "",
        Port: 0,
        SourceAddr: "",
        Description: "",

        ConnectionNumber: 0,
        UploadSpeed: [0],
        DownloadSpeed: [0],

        Capacity: CLIENT_CAPACITY,

        ChartConfig: {
            // The type of chart we want to create
            type: 'line',
        
            // The data for our dataset
            data: {
                labels: new Array(CLIENT_CAPACITY),
                datasets: 
                [
                    {
                        label: 'up: ',
                        backgroundColor: 'rgb(0, 255, 128)',
                        borderColor: 'rgb(0, 255, 128)',
                        data: this.UploadSpeed,
                        fill: false,
                        borderWidth: 1,
                        pointRadius: 1
                    },
                    {
                        label: 'down: ',
                        backgroundColor: 'rgb(0, 128, 255)',
                        borderColor: 'rgb(0, 128, 255)',
                        data: this.DownloadSpeed,
                        fill: false,
                        borderWidth: 1,
                        pointRadius: 1
                    }
                ]
            },
        
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: {
                    duration: 0
                }
            }
        },

        SetDiv: function(divid){
            this.DivId = divid
        },

        Update: function(c){
            this.ClientId = c.ClientId
            this.ClientAddr = c.ClientAddr
            this.Port = c.Port
            this.SourceAddr = c.SourceAddr
            this.Description = c.Description

            this.ConnectionNumber = c.ConnectionNumber
            var us = c.UploadSpeed, ds = c.DownloadSpeed
            this.UploadSpeed.push(us)
            this.DownloadSpeed.push(ds)

            while(this.UploadSpeed.length > this.Capacity){
                this.UploadSpeed.shift()
            }

            while(this.DownloadSpeed.length > this.Capacity){
                this.DownloadSpeed.shift()
            }

            //var dom = new DOMParser().parseFromString(this.HTML(), "text/html")
            //document.getElementById(this.DivId).appendChild(a.childNodes[0])  
            //var divDom = document.getElementById(this.DivId).childNodes[0] = dom          
            document.getElementById(this.DivId).innerHTML = this.HTML()
            this.FreshChart()
        },

        FreshChart: function(){
            var ctx = document.getElementById('canvas_' + this.ClientId).getContext('2d');
            var chart = new Chart(ctx, this.ChartConfig)

            this.ChartConfig.data.datasets[0].data = this.UploadSpeed
            this.ChartConfig.data.datasets[1].data = this.DownloadSpeed 
            uploadSpeed = this.UploadSpeed[this.UploadSpeed.length - 1]
            downloadSpeed = this.DownloadSpeed[this.DownloadSpeed.length - 1]
            this.ChartConfig.data.datasets[0].label = "up: " + SpeedToReadable(uploadSpeed)
            this.ChartConfig.data.datasets[1].label = "down: " + SpeedToReadable(downloadSpeed)
        },

        HTML: function(){
            var res = 
            '<div class="row">' +
                '<div class="col-sm-4">' +
                    '<div class="row">' +
                        '<div class="col-sm-4"><h6>ClientId:</h6></div>' +
                        '<div class="col-sm-8">' + this.ClientId + '</div>' +
                    '</div>' +
                    
                    '<div class="row">' +
                        '<div class="col-sm-4"><h6>ClientAddr:</h6></div>' +
                        '<div class="col-sm-8">' + this.ClientAddr + '</div>' + 
                    '</div>' +

                    '<div class="row">' +
                        '<div class="col-sm-4"><h6>SourceAddr:</h6></div>' +
                        '<div class="col-sm-8">' + this.SourceAddr + '</div>' + 
                    '</div>' +

                    '<div class="row">' +
                        '<div class="col-sm-4"><h6>PortTo:</h6></div>' +
                        '<div class="col-sm-8">' + this.Port + '</div>' + 
                    '</div>' +

                    '<div class="row">' +
                        '<div class="col-sm-4"><h6>Description:</h6></div>' +
                        '<div class="col-sm-8">' + this.Description + '</div>' + 
                    '</div>' +

                    '<div class="row">' +
                        '<div class="col-sm-4"><h6>Connections:</h6></div>' +
                        '<div class="col-sm-8">' + this.ConnectionNumber + '</div>' + 
                    '</div>' +

                '</div>' + 

                '<div class="col-sm-8">' +
                    '<canvas id="' + 'canvas_' + this.ClientId + '"></canvas>' +
                '</div>' +
            '</div>' 

            return res
        }
    }
}