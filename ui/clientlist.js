function ClientList(divid){
    return {
        DivId: divid,
        Clients: [],

        SetDiv: function(divid){
            this.DivId = divid
        },

        Update: function(clients){
            var cmap = {}, cmapnew = {}
            for(var i=0; i<clients.length; i++){
                cmapnew[clients[i].ClientId] = i
            }
            leftClients = []
            for(var i=0; i<this.Clients.length; i++){
                if(this.Clients[i].ClientId in cmapnew){
                    leftClients.push(this.Clients[i])
                }
            }
            this.Clients = leftClients
            for(var i=0; i<this.Clients.length; i++){
                cmap[this.Clients[i].ClientId] = i
            }

            for(var i=0; i<clients.length; i++){
                var c = clients[i]
                if(!(c.ClientId in cmap)){
                    this.Clients.push(Client(""))
                    var idx = this.Clients.length - 1
                    cmap[c.ClientId] = idx
                }

                var idx = cmap[c.ClientId]
                this.Clients[idx].SetDiv("itemdiv_" + idx)
            }


            document.getElementById(this.DivId).innerHTML = this.HTML()

            for(var i=0; i<clients.length; i++){
                var c = clients[i]
                var ci = cmap[c.ClientId]
                this.Clients[ci].Update(c)
            }
        },

        HTML: function(){
            var res = 
                '<div class="alert alert-warning" role="alert">' + 
                    'No clients connected.' + 
                '</div>'
            if(this.Clients.length > 0){
                res = ''
            }

            for(var i=0; i<this.Clients.length; i++){
                res += "<div style='margin:10px; padding: 10px;' class='card bg-light' id='" + "itemdiv_" + i + "'></div>"
            }

            return res
        }
    }
}