function Server(divid){
    uploadSpeed = []
    downloadSpeed = []
    for(var i=0; i<RECORD_CAPACITY; i++){
        uploadSpeed.push(-1)
        downloadSpeed.push(-1)
    }

    return {
        DivId: divid,
        ServerAddr: "",

        TCPClientNumber: 0,
        UDPClientNumber: 0,
        TCPConnectionNumber: 0,
        UDPConnectionNumber: 0,
        ForwardNumber: 0,
        ReverseNumber: 0,

        UploadSpeed: uploadSpeed,
        DownloadSpeed: downloadSpeed,

        Capacity: RECORD_CAPACITY,

        ClientChartConfig: {
			type: 'pie',
			data: {
				datasets: [{
					data: [
                        0,
                        0
					],
					backgroundColor: [
						'rgb(0, 255, 0, 0.5)',
						'rgb(0, 0, 255, 0.5)'
					],
					label: 'Client Number'
				}],
				labels: [
					'tcp',
					'udp',
				]
			},
			options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: {
                    duration: 0
                },

                layout: {
                    padding: {
                        left: 0,
                        right: 0,
                        top: 0,
                        bottom: 0
                    }
                },

                title: {
                    display: true,
                    text: 'tcp/udp clients',
                    position: 'left'
                }
			}
        },

        ConnectionChartConfig: {
			type: 'pie',
			data: {
				datasets: [{
					data: [
                        0,
                        0
					],
					backgroundColor: [
						'rgb(0, 255, 0, 0.5)',
						'rgb(0, 0, 255, 0.5)'
					],
					label: 'Client Number'
				}],
				labels: [
					'tcp',
					'udp',
				]
			},
			options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: {
                    duration: 0
                },

                layout: {
                    padding: {
                        left: 0,
                        right: 0,
                        top: 0,
                        bottom: 0
                    }
                },

                title: {
                    display: true,
                    text: 'connections',
                    position: 'left'
                }
			}
        },
        
        DirectionChartConfig: {
			type: 'pie',
			data: {
				datasets: [{
					data: [
                        0,
                        0
					],
					backgroundColor: [
						'rgb(0, 255, 0, 0.5)',
						'rgb(0, 0, 255, 0.5)'
					],
					label: 'Direction'
				}],
				labels: [
					'forward',
					'reverse',
				]
			},
			options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: {
                    duration: 0
                },

                layout: {
                    padding: {
                        left: 0,
                        right: 0,
                        top: 0,
                        bottom: 0
                    }
                },

                title: {
                    display: true,
                    text: 'f/r clients',
                    position: 'left'
                }
			}
		},

        SpeedChartConfig: {
            // The type of chart we want to create
            type: 'line',
        
            // The data for our dataset
            data: {
                labels: new Array(RECORD_CAPACITY),
                datasets: 
                [
                    {
                        label: 'up: ',
                        backgroundColor: 'rgb(0, 255, 128, 0.3)',
                        borderColor: 'rgb(0, 255, 128)',
                        data: this.UploadSpeed,
                        fill: true,
                        borderWidth: 2,
                        pointRadius: 0
                    },
                    {
                        label: 'down: ',
                        backgroundColor: 'rgb(0, 128, 255, 0.3)',
                        borderColor: 'rgb(0, 128, 255)',
                        data: this.DownloadSpeed,
                        fill: true,
                        borderWidth: 2,
                        pointRadius: 0
                    }
                ]
            },
        
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: {
                    duration: 0
                },

                layout: {
                    padding: {
                        left: 0,
                        right: 0,
                        top: 0,
                        bottom: 0
                    }
                },

                scales: {
					yAxes: [{
						display: true,
						scaleLabel: {
							display: true,
							labelString: 'B/s'
                        },
                        ticks: {
                            min: 0
                        }
                    }]
				}
            }
        },

        SetDiv: function(divid){
            this.DivId = divid
        },

        Update: function(c){
            this.ServerAddr = c.ServerAddr
            var uploadSpeed = 0, downloadSpeed = 0
            var tcpConnNumber = 0, udpConnNumber = 0
            var tcpClientNumber = 0, udpClientNumber = 0
            var forwardNumber = 0, reverseNumber = 0
            for(var i=0; i<c.Clients.length; i++){
                uploadSpeed += c.Clients[i].UploadSpeed
                downloadSpeed += c.Clients[i].DownloadSpeed

                if(c.Clients[i].Protocol == "tcp"){
                    tcpConnNumber += c.Clients[i].ConnectionNumber
                    tcpClientNumber += 1
                }
                if(c.Clients[i].Protocol == "udp"){
                    udpConnNumber += c.Clients[i].ConnectionNumber
                    udpClientNumber += 1
                }

                if(c.Clients[i].Direction == "forward"){
                    forwardNumber += 1
                }

                if(c.Clients[i].Direction == "reverse"){
                    reverseNumber += 1
                }
            }

            this.TCPClientNumber = tcpClientNumber
            this.UDPClientNumber = udpClientNumber
            this.TCPConnectionNumber = tcpConnNumber
            this.UDPConnectionNumber = udpConnNumber
            this.ForwardNumber = forwardNumber
            this.ReverseNumber = reverseNumber
           
            this.UploadSpeed.push(uploadSpeed)
            this.DownloadSpeed.push(downloadSpeed)

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
            var speedCtx = document.getElementById('canvas_speed').getContext('2d');
            var speedChart = new Chart(speedCtx, this.SpeedChartConfig)

            this.SpeedChartConfig.data.datasets[0].data = this.UploadSpeed
            this.SpeedChartConfig.data.datasets[1].data = this.DownloadSpeed 
            uploadSpeed = this.UploadSpeed[this.UploadSpeed.length - 1]
            downloadSpeed = this.DownloadSpeed[this.DownloadSpeed.length - 1]
            this.SpeedChartConfig.data.datasets[0].label = "up: " + SpeedToReadable(uploadSpeed)
            this.SpeedChartConfig.data.datasets[1].label = "down: " + SpeedToReadable(downloadSpeed)

            var clientCtx = document.getElementById('canvas_client').getContext('2d');
            var clientChart = new Chart(clientCtx, this.ClientChartConfig)
            this.ClientChartConfig.data.datasets[0].data = [this.TCPClientNumber, this.UDPClientNumber]
            //this.ClientChartConfig.data.labels = ["tcp: " + this.TCPClientNumber, "udp: " + this.UDPClientNumber]

            var connectionCtx = document.getElementById('canvas_connection').getContext('2d');
            var connectionChart = new Chart(connectionCtx, this.ConnectionChartConfig)
            this.ConnectionChartConfig.data.datasets[0].data = [this.TCPConnectionNumber, this.UDPConnectionNumber]
            //this.ClientChartConfig.data.labels = ["tcp: " + this.TCPClientNumber, "udp: " + this.UDPClientNumber]

            var directionCtx = document.getElementById('canvas_direction').getContext('2d');
            var directionChart = new Chart(directionCtx, this.DirectionChartConfig)
            this.DirectionChartConfig.data.datasets[0].data = [this.ForwardNumber, this.ReverseNumber]
            //this.ClientChartConfig.data.labels = ["tcp: " + this.TCPClientNumber, "udp: " + this.UDPClientNumber]
        },

        HTML: function(){
            var res = 
            '<div class="row">' +
                '<div class="col-sm-2">' +
                    '<canvas id="canvas_client"></canvas>' +
                '</div>' +

                '<div class="col-sm-2">' +
                    '<canvas id="canvas_direction"></canvas>' +
                '</div>' +

                '<div class="col-sm-2">' +
                    '<canvas id="canvas_connection"></canvas>' +
                '</div>' +

                '<div class="col-sm-6">' +
                    '<canvas id="canvas_speed"></canvas>' +
                '</div>' +

            '</div>' 

            return res
        }
    }
}