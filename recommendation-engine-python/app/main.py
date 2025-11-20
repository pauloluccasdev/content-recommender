from fastapi import FastAPI
from app.api.routes import recommendations
from app.core.config import settings

# Criar instância do FastAPI
app = FastAPI(
    title=settings.app_name,
    version=settings.version,
    description="API de sistema de recomendação de conteúdos usando Pandas + Scikit-Learn"
)

# Incluir rotas
app.include_router(recommendations.router)

@app.get("/")
def read_root():
    """Endpoint raiz da API"""
    return {
        "mensagem": "Bem-vindo à API de Recomendação de Conteúdos!",
        "version": settings.version,
        "docs": "/docs"
    }

@app.get("/health")
def health_check():
    """Endpoint de health check"""
    return {"status": "healthy"}