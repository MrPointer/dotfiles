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
