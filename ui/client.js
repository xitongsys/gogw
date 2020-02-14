MSG_FRESH_CHART = "fresh_chart_"

function Client(divid){
    return {
        DivId: divid,
        ClientId: "0",
        ClientAddr: "",
        Port: 0,
        SourceAddr: "",
        Description: "",

        ConnectionNumber: 0,
        UploadSpeed: [0,0,0],
        DownloadSpeed: [0,0,0],

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
                        label: 'Upload Speed',
                        backgroundColor: 'rgb(0, 255, 128)',
                        borderColor: 'rgb(0, 255, 128)',
                        data: this.UploadSpeed,
                        fill: false
                    },
                    {
                        label: 'Download Speed',
                        backgroundColor: 'rgb(0, 128, 255)',
                        borderColor: 'rgb(0, 128, 255)',
                        data: this.DownloadSpeed,
                        fill: false
                    }
                ]
            },
        
            // Configuration options go here
            options: {
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

            while(this.UploadSpeed.length > this.capacity){
                this.UploadSpeed.shift()
                console.log(this.UploadSpeed)
            }

            while(this.DownloadSpeed.length > this.capacity){
                this.DownloadSpeed.shift()
            }

            document.getElementById(this.DivId).innerHTML = this.HTML()
            this.FreshChart()
        },

        FreshChart: function(){
            var ctx = document.getElementById('canvas_' + this.ClientId).getContext('2d');
            var chart = new Chart(ctx, this.ChartConfig)
            //this.FirstFresh = false
            this.ChartConfig.data.datasets[0].data = this.UploadSpeed
            this.ChartConfig.data.datasets[1].data = this.DownloadSpeed 

        },

        HTML: function(){
            var res = 
            '<div class="row">' +
                '<div class="col-sm-6">' +
                    '<div class="row">' +
                        '<div class="col-sm-4">ClientId </div>' +
                        '<div class="col-sm-8">' + this.ClientId + '</div>' +
                    '</div>' +
                    
                    '<div class="row">' +
                        '<div class="col-sm-4">ClientAddr </div>' +
                        '<div class="col-sm-8">' + this.ClientAddr + '</div>' + 
                    '</div>' +

                    '<div class="row">' +
                        '<div class="col-sm-4">SourceAddr </div>' +
                        '<div class="col-sm-8">' + this.SourceAddr + '</div>' + 
                    '</div>' +

                    '<div class="row">' +
                        '<div class="col-sm-4">PortTo </div>' +
                        '<div class="col-sm-8">' + this.Port + '</div>' + 
                    '</div>' +

                    '<div class="row">' +
                        '<div class="col-sm-4">Description </div>' +
                        '<div class="col-sm-8">' + this.Description + '</div>' + 
                    '</div>' +

                    '<div class="row">' +
                        '<div class="col-sm-4">' + this.ConnectionNumber + '</div>' +
                        '<div class="col-sm-4">' + this.UploadSpeed[this.UploadSpeed.length-1] + '</div>' + 
                        '<div class="col-sm-4">' + this.DownloadSpeed[this.DownloadSpeed.length-1] + '</div>' + 
                    '</div>' +
                '</div>' + 

                '<div class="col-sm-6">' +
                    '<canvas id="' + 'canvas_' + this.ClientId + '"></canvas>' +
                '</div>' +
            '</div>' 

            return res
        }
    }
}