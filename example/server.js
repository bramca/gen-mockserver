
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
	"/auth/register": "/auth-register",
	"/auth/login": "/auth-login",
	"/cart/items": "/cart-items",
	"/checkout": "/checkout",
	"/addresses": "/addresses",
	"/orders": "/orders",
	"/orders/:orderId": "/orders/:orderId",
	"/addresses": "/addresses",
	"/products": "/products",
	"/products?category=": "/products",
	"/products?search=": "/products",
	"/products?min_price=": "/products",
	"/products?max_price=": "/products",
	"/products/:id": "/products/:id",
	"/cart": "/cart"
}));

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

server.post('/checkout', (req, res) => {
	console.log(`POST /checkout with body ${JSON.stringify(req.body)}`);
	statusCode = 201;
	responseBody = {
		"created_at": "1909-03-05 07:48:52.266464515 +0000 UTC",
		"id": "3095d2d1-0e92-47ab-8c12-f99b4f4973e3",
		"items": [
				{
						"product_id": "ab771e86-1957-44b0-9f24-7dc71130dab0",
						"quantity": 43
					}
		],
		"status": "confirmed",
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

server.get('/orders', (req, res) => {
	console.log(`GET /orders`);
	statusCode = 200;
	responseBody = [
		{
				"created_at": "1924-04-16 03:00:39.061914978 +0000 UTC",
				"id": "04f7d791-3af3-4ddf-9dec-442949b34c59",
				"items": [],
				"status": "confirmed",
				"total_amount": null
			}
];
	
	res.status(statusCode).json(responseBody);
});

server.get('/orders/:orderId', (req, res) => {
	console.log(`GET /orders/${req.params.orderId}`);
	statusCode = 200;
	responseBody = {
		"created_at": "1962-03-24 07:00:16.941815842 +0000 UTC",
		"id": "2c81fb08-d4bf-4737-94cf-c6a1a2727666",
		"items": [
				{
						"product_id": "c7d721db-1e45-404f-b592-b63e2ceafb21",
						"quantity": 51
					}
		],
		"status": "shipped",
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
				"country": "qOZss",
				"line1": "hkJHjiFYGPbYNa",
				"line2": "",
				"postal_code": "eJI",
				"state": "nUbOaD"
			}
];
	
	res.status(statusCode).json(responseBody);
});

server.get('/products', (req, res) => {
	console.log(`GET /products`);
	statusCode = 200;
	responseBody = [
		{
				"category": "gtQdDNCfcPI",
				"created_at": "1905-01-13 17:30:17.849755166 +0000 UTC",
				"description": "O",
				"id": "d75d6f96-f9c4-4361-b279-30c861733362",
				"image_url": "tsnoUoBYRctVKb",
				"name": "nAnu",
				"price": null,
				"stock": 45,
				"updated_at": "1984-12-29 23:07:05.522977938 +0000 UTC"
			}
];
	
	res.status(statusCode).json(responseBody);
});

server.get('/products/:id', (req, res) => {
	console.log(`GET /products/${req.params.id}`);
	statusCode = 200;
	responseBody = {
		"category": "QbvLjzXd",
		"created_at": "2016-06-24 09:45:25.159526303 +0000 UTC",
		"description": "jGstYdKzQ",
		"id": "934fe19a-f2a3-459b-88dc-fa079dfa6c44",
		"image_url": "wALExBvOXWRuka",
		"name": "zfb",
		"price": null,
		"stock": 23,
		"updated_at": "2007-11-20 23:43:43.51270311 +0000 UTC"
	};
	
	res.status(statusCode).json(responseBody);
});

server.get('/cart', (req, res) => {
	console.log(`GET /cart`);
	statusCode = 200;
	responseBody = [
		{
				"product_id": "a44ed77e-ad72-46f5-bb95-9cba82ece0b8",
				"quantity": 75
			}
];
	
	res.status(statusCode).json(responseBody);
});

server.use(middlewares);
server.use(router);
server.listen(port);
