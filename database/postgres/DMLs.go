package postgres

const AddClient = `insert into clients(name, surname, login, password, age, gender, phone , status)
	values($1, $2, $3, $4, $5, $6, $7, $8)`

const AddAccount = `insert into accounts( client_id, account_number, balance, status, card_number )
	values( $1, $2, $3, $4, $5 )`

const updateLimit = `update accounts
	set limit_transfer = $limit_transfer, limit_payment = $limit_payment where id = $accountId`

const AddATMs = `insert into atms( name, status )
	values( $1, $2 )`

const changeStatusATM = `update atms
	set status = $status where id = $atmId`

const AddService = `insert into services( name, account_number )
	values( $1, $2 )`

const GetAllClients = `select * from clients`

const GetAllAccounts = `select a.id, a.client_id, a.account_number, a.balance, a.status,
									c.id,c.name,c.surname,c.login,c.password,c.phone,c.status,c.verified_at 
						from accounts a left join clients c on a.client_id = c.id`

const GetAllATMs = `select * from ATMs`

const LoginSQL = `select * from clients where login = ($1)`

const SearchClientByLogin = `select id, surname from clients where login = ($1)`

const SearchAccountByID = `select id, client_id, account_number, balance, status, card_number from accounts where status = true and client_id = ($1)`

const GetAllServices = `select id, name from services`

