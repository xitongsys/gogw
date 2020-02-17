function SpeedToReadable(s){
    if(s <= 1024){
        return s + " B/s"
    }else if(s < 1024 * 1024){
        return (s/1024.0).toFixed(2) + " KB/s"
    }else{
        return (s/1024.0/1024.0).toFixed(2) + " MB/s"
    }
}
