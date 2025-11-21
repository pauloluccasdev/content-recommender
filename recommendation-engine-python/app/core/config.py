from pydantic_settings import BaseSettings
from typing import Optional

class Settings(BaseSettings):
    """Configurações da aplicação"""
    app_name: str = "Content Recommender Engine API"
    version: str = "1.0.0"
    default_top_n: int = 10
    
    # Configurações de banco de dados
    db_host: str = "db"
    db_port: int = 3306
    db_user: str = "root"
    db_password: str = "vertrigo"
    db_name: str = "content-recommender"
    database_url: Optional[str] = None
    
    # Modo de operação: 'real' usa banco de dados, 'simulated' usa dados simulados
    data_mode: str = "real"  # 'real' ou 'simulated'
    
    # Configurações para dados simulados (fallback)
    simulated_users: int = 100
    simulated_contents: int = 50

    class Config:
        env_file = ".env"
        case_sensitive = False

    def get_database_url(self) -> str:
        """Gera a URL de conexão do banco de dados"""
        if self.database_url:
            return self.database_url
        
        return f"mysql+pymysql://{self.db_user}:{self.db_password}@{self.db_host}:{self.db_port}/{self.db_name}"

settings = Settings()