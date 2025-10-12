" Load vim-plug
call plug#begin(stdpath('data') . '/plugged')

" Shorthand notation; fetches https://github.com/junegunn/vim-easy-align
Plug 'junegunn/vim-easy-align'
Plug 'neoclide/coc.nvim', {'branch': 'release'}
Plug 'joshdick/onedark.vim'
Plug 'ayu-theme/ayu-vim'
Plug 'sheerun/vim-polyglot'
Plug 'vim-airline/vim-airline'
Plug 'vim-airline/vim-airline-themes'

" Initialize plugin system
call plug#end()

syntax on
let ayucolor="dark"
colorscheme ayu
let g:airline_theme='ayu_dark'

" Configure coc plugins
let g:coc_global_extensions = [
    \'coc-clangd',
    \'coc-cmake',
    \'coc-git',
    \'coc-go',
    \'coc-highlight',
    \'coc-json',
    \'coc-markdownlint',
    \'coc-pyright',
    \'coc-rust-analyzer',
    \'coc-sh',
    \'coc-yaml',
\]

filetype plugin indent on
" show existing tab with 4 spaces width
set tabstop=4
" when indenting with '>', use 4 spaces width
set shiftwidth=4
" On pressing tab, insert 4 spaces
set expandtab
" Show line numbers
set number

" CoC completion
inoremap <expr> <cr> coc#pum#visible() ? coc#pum#confirm() : "\<CR>"

" Tab navigation
nnoremap <A-Left> :tabprevious<CR>
nnoremap <A-Right> :tabnext<CR>

" Search and replace helpers
vnoremap <C-r> "hy:%s/<C-r>h//gc<left><left><left>
vnoremap <C-f> y<ESC>/<c-r>"<CR>

" Essential Alt/Option key mappings for word navigation
" These work across most terminals on macOS
inoremap <A-Left>  <C-o>b
inoremap <A-Right> <C-o>w
inoremap <M-BS>    <C-o>db
inoremap <M-Del>   <C-o>de

" Add fallback escape sequences for Alt+Left (common terminal sequences)
execute "inoremap \e[1;3D <C-o>b"
execute "inoremap \e[D <C-o>b"
execute "inoremap \eb <C-o>b"
execute "nnoremap \e[1;3D b"
execute "nnoremap \e[D b"
execute "nnoremap \eb b"
execute "vnoremap \e[1;3D b"
execute "vnoremap \e[D b"
execute "vnoremap \eb b"

" Add fallback escape sequences for Alt+Right (common terminal sequences)
execute "inoremap \e[1;3C <C-o>w"
execute "inoremap \e[C <C-o>w"
execute "inoremap \ef <C-o>w"
execute "nnoremap \e[1;3C w"
execute "nnoremap \e[C w"
execute "nnoremap \ef w"
execute "vnoremap \e[1;3C w"
execute "vnoremap \e[C w"
execute "vnoremap \ef w"

" Add fallback escape sequences for Alt+Delete (common terminal sequences)
execute "inoremap \e[3;3~ <C-o>de"
execute "inoremap \ed <C-o>de"

" Normal mode equivalents
nnoremap <A-Left>  b
nnoremap <A-Right> w

" Visual mode equivalents
vnoremap <A-Left>  b
vnoremap <A-Right> w

" Home/End key support (most common sequences)
inoremap <Home> <C-o>0
inoremap <End>  <C-o>$
inoremap <C-A>  <C-o>0
inoremap <C-E>  <C-o>$

nnoremap <Home> 0
nnoremap <End>  $

vnoremap <Home> 0
vnoremap <End>  $
