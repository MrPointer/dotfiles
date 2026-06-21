document.addEventListener('click',e=>{
  const choice=e.target.closest('[data-theme-choice]');
  if(choice){
    const theme=choice.dataset.themeChoice;
    if(theme==='system')delete document.documentElement.dataset.theme;
    else document.documentElement.dataset.theme=theme;
    document.querySelectorAll('[data-theme-choice]').forEach(b=>b.setAttribute('aria-pressed',String(b===choice)));
    return;
  }
  const t=e.target.closest('.tab');
  if(t){const card=t.dataset.card,name=t.dataset.tab;
    document.querySelectorAll(`.tab[data-card="${card}"]`).forEach(b=>b.classList.toggle('active',b===t));
    document.querySelectorAll(`.pane[data-card="${card}"]`).forEach(p=>p.classList.toggle('active',p.dataset.pane===name));
    return;}
  const node=e.target.closest('.node');
  if(node){const id=node.dataset.target,el=document.getElementById(id);
    if(el){el.open=true;el.scrollIntoView({behavior:'smooth',block:'start'});}}
});
const ex=document.getElementById('expand'),co=document.getElementById('collapse');
if(ex)ex.onclick=()=>document.querySelectorAll('details.card').forEach(d=>d.open=true);
if(co)co.onclick=()=>document.querySelectorAll('details.card').forEach(d=>d.open=false);

function openPlanCard(id){
  const el=document.getElementById(id);
  if(el){el.open=true;el.scrollIntoView({behavior:'smooth',block:'start'});}
}

document.querySelectorAll('[data-graph-surface]').forEach(surface=>{
  const dataEl=surface.querySelector('[data-plan-graph]'),canvas=surface.querySelector('[data-graph-canvas]');
  if(!dataEl||!canvas||typeof dagre==='undefined')return;
  let data;
  try{data=JSON.parse(dataEl.textContent||'{}');}catch{return;}
  if(!data.nodes||!data.nodes.length)return;

  const svgNS='http://www.w3.org/2000/svg';
  const svg=document.createElementNS(svgNS,'svg');
  const inner=document.createElementNS(svgNS,'g');
  const defs=document.createElementNS(svgNS,'defs');
  const marker=document.createElementNS(svgNS,'marker');
  marker.setAttribute('id','graph-arrow');marker.setAttribute('viewBox','0 0 10 10');
  marker.setAttribute('refX','9');marker.setAttribute('refY','5');marker.setAttribute('markerWidth','7');marker.setAttribute('markerHeight','7');marker.setAttribute('orient','auto-start-reverse');
  const arrow=document.createElementNS(svgNS,'path');
  arrow.setAttribute('d','M 0 0 L 10 5 L 0 10 z');marker.appendChild(arrow);defs.appendChild(marker);svg.appendChild(defs);svg.appendChild(inner);canvas.appendChild(svg);

  const g=new dagre.graphlib.Graph().setGraph({rankdir:data.direction==='TD'?'TB':'LR',nodesep:42,ranksep:72,marginx:24,marginy:24}).setDefaultEdgeLabel(()=>({}));
  data.nodes.forEach(n=>g.setNode(n.id,{label:n.label,subplan:n.subplan,hue:n.hue,width:176,height:58}));
  data.edges.forEach(e=>g.setEdge(e.source,e.target));
  dagre.layout(g);

  function path(points){return points&&points.length?points.map((p,i)=>(i?'L':'M')+p.x+' '+p.y).join(' '):'';}
  g.edges().forEach(e=>{
    const p=document.createElementNS(svgNS,'path');
    p.setAttribute('class','graph-edge');p.setAttribute('d',path(g.edge(e).points));p.setAttribute('marker-end','url(#graph-arrow)');inner.appendChild(p);
  });
  g.nodes().forEach(id=>{
    const n=g.node(id),group=document.createElementNS(svgNS,'g');
    group.setAttribute('class','graph-node');group.setAttribute('tabindex','0');group.setAttribute('role','button');group.setAttribute('data-target','sp-'+n.subplan);
    group.setAttribute('transform','translate('+n.x+' '+n.y+')');
    if(n.hue)group.style.setProperty('--hue',n.hue);
    const rect=document.createElementNS(svgNS,'rect');
    rect.setAttribute('x',-n.width/2);rect.setAttribute('y',-n.height/2);rect.setAttribute('width',n.width);rect.setAttribute('height',n.height);rect.setAttribute('rx','10');
    const text=document.createElementNS(svgNS,'text');
    text.setAttribute('text-anchor','middle');text.setAttribute('dominant-baseline','middle');text.textContent=n.label;
    group.appendChild(rect);group.appendChild(text);inner.appendChild(group);
  });

  const graph=g.graph(),pad=20;
  let scale=1,drag=null,suppressClick=false;
  const base={x:-pad,y:-pad,w:(graph.width||0)+pad*2,h:(graph.height||0)+pad*2};
  const view={...base};
  function apply(){svg.setAttribute('viewBox',[view.x,view.y,view.w,view.h].join(' '));}
  function setScale(next){
    const cx=view.x+view.w/2,cy=view.y+view.h/2;
    scale=next;view.w=base.w/scale;view.h=base.h/scale;view.x=cx-view.w/2;view.y=cy-view.h/2;apply();
  }
  apply();
  surface.querySelectorAll('[data-graph-action]').forEach(b=>b.addEventListener('click',()=>{
    const a=b.dataset.graphAction;
    if(a==='fit'){scale=1;Object.assign(view,base);apply();}
    if(a==='zoom-in')setScale(Math.min(2.5,scale*1.2));
    if(a==='zoom-out')setScale(Math.max(.5,scale/1.2));
  }));
  svg.addEventListener('pointerdown',e=>{
    svg.setPointerCapture(e.pointerId);svg.classList.add('is-panning');
    drag={id:e.pointerId,x:e.clientX,y:e.clientY,view:{...view},moved:false};
  });
  svg.addEventListener('pointermove',e=>{
    if(!drag||drag.id!==e.pointerId)return;
    const dx=e.clientX-drag.x,dy=e.clientY-drag.y;
    if(Math.abs(dx)+Math.abs(dy)>3){drag.moved=true;suppressClick=true;}
    const sx=view.w/Math.max(1,svg.clientWidth),sy=view.h/Math.max(1,svg.clientHeight);
    view.x=drag.view.x-dx*sx;view.y=drag.view.y-dy*sy;apply();
  });
  function stopPan(e){if(drag&&drag.id===e.pointerId){svg.classList.remove('is-panning');drag=null;setTimeout(()=>{suppressClick=false;},0);}}
  svg.addEventListener('pointerup',stopPan);svg.addEventListener('pointercancel',stopPan);
  surface.addEventListener('click',e=>{if(suppressClick)return;const node=e.target.closest('.graph-node');if(node)openPlanCard(node.dataset.target);});
  surface.addEventListener('keydown',e=>{const node=e.target.closest('.graph-node');if(node&&(e.key==='Enter'||e.key===' ')){e.preventDefault();openPlanCard(node.dataset.target);}});
});
