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

            //var dom = new DOMParser().parseFromString(this.HTML(), "text/html")
            //document.getElementById(this.DivId).appendChild(dom.childNodes[0])
            var divDom = document.getElementById(this.DivId)
            if(divDom.childElementCount < this.Clients.length){
                for(var i=divDom.childElementCount; i<this.Clients.length; i++){
                    var domItem = new DOMParser().parseFromString(this.HTMLItem(i), "text/html")
                    divDom.appendChild(domItem.childNodes[0])
                }
            }else if(divDom.childElementCount > this.Clients.length){
                for(var i=divDom.childElementCount-1; i>=this.Clients.length; i--){
                    divDom.childNodes[i].remove()
                }
            }

            for(var i=0; i<clients.length; i++){
                var c = clients[i]
                var ci = cmap[c.ClientId]
                this.Clients[ci].Update(c)
            }
        },

        HTMLItem: function(i){
            return "<div style='margin:10px; padding: 10px;' class='card bg-light' id='" + "itemdiv_" + i + "'></div>"
        },

        HTMLNoItem: function(){
            var res = 
                '<div class="alert alert-warning" role="alert">' + 
                    'No clients connected.' + 
                '</div>'
            return res
        },

        HTML: function(){
            var res = this.HTMLNoItem
            if(this.Clients.length > 0){
                res = ''
            }

            for(var i=0; i<this.Clients.length; i++){
                res += this.HTMLItem(i)
            }
            return res
        }
    }
}