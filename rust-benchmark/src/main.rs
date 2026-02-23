mod entities;

use axum::{
    Json, Router,
    extract::State,
    http::StatusCode,
    response::Html,
    routing::{get, post},
};
use bcrypt::{DEFAULT_COST, hash, verify};
use entities::user;
use sea_orm::{
    ActiveModelTrait, ActiveValue::Set, ColumnTrait, ConnectionTrait, Database,
    DatabaseConnection, DbBackend, EntityTrait, QueryFilter, Statement,
};
use serde::Deserialize;
use tower_http::trace::{DefaultMakeSpan, DefaultOnRequest, DefaultOnResponse, TraceLayer};
use tracing::Level;

#[derive(Clone)]
struct AppState {
    db: DatabaseConnection,
}

#[derive(Deserialize)]
struct AuthPayload {
    email: String,
    password: String,
}

async fn handler() -> Html<&'static str> {
    Html("<h1>Hello, World!</h1>")
}

async fn register(
    State(state): State<AppState>,
    Json(payload): Json<AuthPayload>,
) -> Result<Json<&'static str>, StatusCode> {
    tracing::info!("register : {}", payload.email);

    let hashed = hash(&payload.password, DEFAULT_COST)
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let new_user = user::ActiveModel {
        email: Set(payload.email.clone()),
        password: Set(hashed),
        ..Default::default()
    };

    new_user
        .insert(&state.db)
        .await
        .map_err(|_| StatusCode::CONFLICT)?;

    Ok(Json("registered"))
}

async fn login(
    State(state): State<AppState>,
    Json(payload): Json<AuthPayload>,
) -> Result<Json<&'static str>, StatusCode> {
    tracing::info!("login attempt : {}", payload.email);

    let user = user::Entity::find()
        .filter(user::Column::Email.eq(&payload.email))
        .one(&state.db)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    match user {
        Some(u) => {
            let valid = verify(&payload.password, &u.password)
                .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
            if valid {
                tracing::info!("login success : {}", payload.email);
                Ok(Json("logged in"))
            } else {
                Err(StatusCode::UNAUTHORIZED)
            }
        }
        None => Err(StatusCode::NOT_FOUND),
    }
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .with_env_filter("info")
        .init();

    dotenvy::dotenv().ok();

    let database_url = std::env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgres://postgres:postgres@localhost:5432/benchmark".to_string());

    let db = Database::connect(&database_url)
        .await
        .expect("Failed to connect to PostgreSQL");

    db.execute(Statement::from_string(
        DbBackend::Postgres,
        "CREATE TABLE IF NOT EXISTS users (
            id       SERIAL PRIMARY KEY,
            email    TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        )"
        .to_string(),
    ))
    .await
    .expect("Failed to create users table");

    tracing::info!("Connected to PostgreSQL");

    let state = AppState { db };

    let app = Router::new()
        .route("/", get(handler))
        .route("/register", post(register))
        .route("/login", post(login))
        .layer(
            TraceLayer::new_for_http()
                .make_span_with(DefaultMakeSpan::new().level(Level::INFO))
                .on_request(DefaultOnRequest::new().level(Level::TRACE))
                .on_response(DefaultOnResponse::new().level(Level::TRACE)),
        )
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
