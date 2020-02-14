function ClientList(divid){
    return {
        DivId: divid,
        Clients: [],

        SetDiv: function(divid){
            this.DivId = divid
        },

        Update: function(clients){
            var cmap = {}
            for(var i=0; i<this.Clients.length; i++){
                cmap[this.Clients[i].ClientId] = i
            }

            for(var i=0; i<clients.length; i++){
                var c = clients[i]
                if(!(c.ClientId in cmap)){
                    this.Clients.push(Client(""))
                    var idx = this.Clients.length - 1
                    cmap[c.ClientId] = idx
                    this.Clients[idx].SetDiv("itemdiv_" + idx)
                }
            }

            document.getElementById(this.DivId).innerHTML = this.HTML()

            for(var i=0; i<clients.length; i++){
                var c = clients[i]
                var ci = cmap[c.ClientId]
                this.Clients[ci].Update(c)
            }
        },

        HTML: function(){
            var res = ""
            for(var i=0; i<this.Clients.length; i++){
                res += "<div class='card bg-light' id='" + "itemdiv_" + i + "'></div>"
            }

            return res
        }
    }
}