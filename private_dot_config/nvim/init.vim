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

inoremap <expr> <cr> coc#pum#visible() ? coc#pum#confirm() : "\<CR>"
nnoremap <A-Left> :tabprevious<CR>
nnoremap <A-Right> :tabnext<CR>
vnoremap <C-r> "hy:%s/<C-r>h//gc<left><left><left>
vnoremap <C-f> y<ESC>/<c-r>"<CR>

