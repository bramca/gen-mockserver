
const fs = require('fs');
const jsonServer = require('json-server');
const db = require('./db.json')
const server = jsonServer.create();
const router = jsonServer.router('./db.json');
const middlewares = jsonServer.defaults();
const port = process.env.PORT || 5000;
const persistentStorage = process.env.STORE || false;

server.use(jsonServer.bodyParser);

function checkWriteToDb() {
	if (persistentStorage) {
		fs.writeFile('./db.json', JSON.stringify(db, undefined, 2), (err) => {
			if (err) throw err;
		});
	}
}

server.use(jsonServer.rewriter({
	"/checkout": "/checkout",
	"/addresses": "/addresses",
	"/auth/register": "/auth-register",
	"/auth/login": "/auth-login",
	"/cart/items": "/cart-items",
	"/products": "/products",
	"/products?category=": "/products",
	"/products?search=": "/products",
	"/products?min_price=": "/products",
	"/products?max_price=": "/products",
	"/products/:id": "/products/:id",
	"/cart": "/cart",
	"/orders": "/orders",
	"/orders/:orderId": "/orders/:orderId",
	"/addresses": "/addresses"
}));

server.post('/checkout', (req, res) => {
	console.log(`POST /checkout with body ${JSON.stringify(req.body)}`);
	statusCode = 201;
	responseBody = {
		"created_at": "",
		"id": "",
		"items": [
				{
						"product_id": "",
						"quantity": 1
					}
		],
		"status": "",
		"total_amount": null
};
	checkWriteToDb();
	res.status(statusCode).json(responseBody);
});

server.post('/addresses', (req, res) => {
	console.log(`POST /addresses with body ${JSON.stringify(req.body)}`);
	statusCode = 201;
	responseBody = undefined;
	checkWriteToDb();
	res.status(statusCode).json(responseBody);
});

server.post('/auth-register', (req, res) => {
	console.log(`POST /auth-register with body ${JSON.stringify(req.body)}`);
	statusCode = 201;
	responseBody = undefined;
	checkWriteToDb();
	res.status(statusCode).json(responseBody);
});

server.post('/auth-login', (req, res) => {
	console.log(`POST /auth-login with body ${JSON.stringify(req.body)}`);
	statusCode = 200;
	responseBody = undefined;
	checkWriteToDb();
	res.status(statusCode).json(responseBody);
});

server.post('/cart-items', (req, res) => {
	console.log(`POST /cart-items with body ${JSON.stringify(req.body)}`);
	statusCode = 200;
	responseBody = undefined;
	checkWriteToDb();
	res.status(statusCode).json(responseBody);
});

server.get('/products', (req, res) => {
	console.log(`GET /products`);
	statusCode = 200;
	responseBody = [
		{
				"category": "",
				"created_at": "",
				"description": "",
				"id": "",
				"image_url": "",
				"name": "",
				"price": null,
				"stock": 0,
				"updated_at": ""
			}
];
	
	res.status(statusCode).json(responseBody);
});

server.get('/products/:id', (req, res) => {
	console.log(`GET /products/${req.params.id}`);
	statusCode = 200;
	responseBody = {
		"category": "",
		"created_at": "",
		"description": "",
		"id": "",
		"image_url": "",
		"name": "",
		"price": null,
		"stock": 0,
		"updated_at": ""
	};
	
	res.status(statusCode).json(responseBody);
});

server.get('/cart', (req, res) => {
	console.log(`GET /cart`);
	statusCode = 200;
	responseBody = [
		{
				"product_id": "",
				"quantity": 1
			}
];
	
	res.status(statusCode).json(responseBody);
});

server.get('/orders', (req, res) => {
	console.log(`GET /orders`);
	statusCode = 200;
	responseBody = [
		{
				"created_at": "",
				"id": "",
				"items": [],
				"status": "",
				"total_amount": null
			}
];
	
	res.status(statusCode).json(responseBody);
});

server.get('/orders/:orderId', (req, res) => {
	console.log(`GET /orders/${req.params.orderId}`);
	statusCode = 200;
	responseBody = {
		"created_at": "",
		"id": "",
		"items": [
				{
						"product_id": "",
						"quantity": 1
					}
		],
		"status": "",
		"total_amount": null
};
	
	res.status(statusCode).json(responseBody);
});

server.get('/addresses', (req, res) => {
	console.log(`GET /addresses`);
	statusCode = 200;
	responseBody = [
		{
				"city": "",
				"country": "",
				"line1": "",
				"line2": "",
				"postal_code": "",
				"state": ""
			}
];
	
	res.status(statusCode).json(responseBody);
});

server.use(middlewares);
server.use(router);
server.listen(port);
