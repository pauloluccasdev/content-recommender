from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    """Configurações da aplicação"""
    app_name: str = "Content Recommender Engine API"
    version: str = "1.0.0"
    default_top_n: int = 10
    simulated_users: int = 100
    simulated_contents: int = 50

    class Config:
        env_file = ".env"

settings = Settings()