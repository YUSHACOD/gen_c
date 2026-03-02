-- .nvim.lua
vim.keymap.set("n", "<C-b>", "<cmd>!go build<cr>")
vim.keymap.set("n", "<C-;>", "<cmd>!go run .<cr>")
vim.keymap.set("n", "<C-x>", ":!go ", { silent = false })
