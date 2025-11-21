from fastapi import FastAPI
from app.api.routes import recommendations
from app.core.config import settings
import logging

# Configurar logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

logger = logging.getLogger(__name__)

# Criar instância do FastAPI
app = FastAPI(
    title=settings.app_name,
    version=settings.version,
    description="API de sistema de recomendação de conteúdos usando Pandas + Scikit-Learn. Integrado com backend-go para usar interações reais dos usuários."
)

# Incluir rotas
app.include_router(recommendations.router)

@app.on_event("startup")
async def startup_event():
    """Evento executado ao iniciar a aplicação"""
    logger.info(f"Iniciando {settings.app_name} v{settings.version}")
    logger.info(f"Modo de dados: {settings.data_mode}")
    
    # Tentar conectar ao banco de dados
    from app.services.database_service import database_service
    if settings.data_mode == "real":
        if database_service.is_connected():
            logger.info("✅ Conectado ao banco de dados")
        else:
            logger.warning("⚠️ Não foi possível conectar ao banco de dados - usando dados simulados")

@app.get("/")
def read_root():
    """Endpoint raiz da API"""
    from app.services.database_service import database_service
    return {
        "mensagem": "Bem-vindo à API de Recomendação de Conteúdos!",
        "version": settings.version,
        "docs": "/docs",
        "data_mode": settings.data_mode,
        "database_connected": database_service.is_connected() if settings.data_mode == "real" else None
    }

@app.get("/health")
def health_check():
    """Endpoint de health check"""
    from app.services.database_service import database_service
    return {
        "status": "healthy",
        "data_mode": settings.data_mode,
        "database_connected": database_service.is_connected() if settings.data_mode == "real" else None
    }