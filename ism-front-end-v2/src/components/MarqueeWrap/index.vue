<template>
  <div class="vue-loop-scroll-box vue-loop-scroll" :class="direction=='left'?'left':''" @mouseenter="stopScroll" @mouseleave="leaveScroll">
    <ul class="vue-loop-scroll-item1">
    <slot></slot>
    </ul>
  <ul class="vue-loop-scroll-item2"></ul>
  </div>
</template>

<script>
	export default{
		 name:'VueLoopScroll',
		 props:{
			 "direction":{default:'up'},
			 "speed":{default:50},
			 "mouseStop":{default:false},
			 "index":{default:0}
		 },
		 mounted(){
			 this.startScroll();
		 },
		 methods:{	
			 //开始滚动
			 startScroll(){
				 this.$nextTick(res=>{
				     var sbox =  document.getElementsByClassName("vue-loop-scroll-box")[this.index];
				     var sitem1 = document.getElementsByClassName("vue-loop-scroll-item1")[this.index];
				     var sitem2 = document.getElementsByClassName("vue-loop-scroll-item2")[this.index];						 
				     sitem2.innerHTML = sitem1.innerHTML;
				     var that = this;
				     this.timer1 = setInterval(function(){
							 if(that.direction=="up"){
									that.scrollUp(sbox,sitem1);
							 }
							 if(that.direction=="left"){
									that.scrollLeft(sbox,sitem1);
							 }
						 },this.speed);
				 })
			 },
							
			 //向上滚动
			 scrollUp(sbox,sitem1){
			    if(sbox.scrollTop>=sitem1.offsetHeight){
			        sbox.scrollTop=0;
			    }else{
			        sbox.scrollTop++;
			    }
			 },	
			 
			 //向左滚动
			 scrollLeft(sbox,sitem1){
			    if(sbox.scrollLeft>=sitem1.offsetWidth){
			        sbox.scrollLeft=0;
			    }else{
			        sbox.scrollLeft++;
			    }
			 },	
			 
			 //鼠标移入停止
			 stopScroll(){
				 if(this.mouseStop){ clearInterval(this.timer1); }
			 },
			 
			 //鼠标移开
			 leaveScroll(){
				 if(this.mouseStop){ this.startScroll(); }
			 }
			 
		 }			
	}
</script>

<style lang="less">
	.vue-loop-scroll{ height:50px; overflow:hidden;
		ul{
			li{ height:50px; overflow:hidden;}
		}
		&.left{ white-space:nowrap;
			ul{ display:inline-block;
				li{ display:inline-block; padding:0 12px;}
			}
		}
	}			
</style>
