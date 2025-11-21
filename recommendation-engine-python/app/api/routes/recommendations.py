"""
Rotas da API para recomendações e interações
"""
from fastapi import APIRouter, HTTPException, BackgroundTasks
from app.models.recommendation import SimpleRecommendationModel
from app.services.dataset_service import DatasetService
from app.services.database_service import database_service
from app.schemas.recommendation import (
    RecommendationRequest,
    RecommendationResponse,
    ContentRecommendation,
    InteractionRequest,
    InteractionResponse
)
from app.core.config import settings
import logging

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/recommendations", tags=["recommendations"])

# Inicializar serviços (singleton)
dataset_service = DatasetService(
    n_users=settings.simulated_users,
    n_contents=settings.simulated_contents
)
recommendation_model = SimpleRecommendationModel(dataset_service)

@router.post("/", response_model=RecommendationResponse)
async def get_recommendations(request: RecommendationRequest):
    """
    Endpoint principal para obter recomendações
    
    Pode usar dois métodos:
    - similarity: Baseado em similaridade entre usuários (collaborative filtering)
    - popularity: Baseado em popularidade dos conteúdos
    """
    try:
        recommendations = recommendation_model.recommend(
            user_id=request.user_id,
            top_n=request.top_n,
            method=request.method
        )
        
        if not recommendations:
            raise HTTPException(
                status_code=404,
                detail=f"Usuário {request.user_id} não encontrado ou sem recomendações disponíveis"
            )
        
        return RecommendationResponse(
            user_id=request.user_id,
            recommendations=[
                ContentRecommendation(**rec) for rec in recommendations
            ],
            method=request.method
        )
    
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Erro interno ao gerar recomendações: {str(e)}"
        )

@router.get("/stats")
async def get_stats():
    """Endpoint para obter estatísticas do dataset"""
    interactions_df = recommendation_model.interactions_df
    contents_df = recommendation_model.contents_df
    
    return {
        "total_users": interactions_df['user_id'].nunique() if interactions_df is not None else 0,
        "total_contents": len(contents_df) if contents_df is not None else 0,
        "total_interactions": len(interactions_df) if interactions_df is not None else 0,
        "avg_interactions_per_user": (
            interactions_df.groupby('user_id').size().mean()
            if interactions_df is not None else 0
        ),
        "data_mode": settings.data_mode,
        "using_real_data": database_service.is_connected() if settings.data_mode == "real" else False
    }

@router.post("/interactions", response_model=InteractionResponse)
async def create_interaction(
    request: InteractionRequest,
    background_tasks: BackgroundTasks
):
    """
    Endpoint para criar uma nova interação de usuário com conteúdo
    
    Quando uma nova interação é criada:
    1. Salva no banco de dados (se em modo real)
    2. Atualiza o modelo de recomendação em background
    """
    try:
        # Salvar interação no banco de dados (se conectado)
        if settings.data_mode == "real" and database_service.is_connected():
            success = database_service.create_interaction(
                user_id=request.user_id,
                content_id=request.content_id,
                interaction_type=request.interaction_type,
                rating=request.rating
            )
            
            if not success:
                raise HTTPException(
                    status_code=500,
                    detail="Erro ao salvar interação no banco de dados"
                )
            
            # Recarregar modelo em background para incluir nova interação
            background_tasks.add_task(reload_recommendation_model)
            logger.info(f"Interação criada e modelo será atualizado em background")
        else:
            # Modo simulado - apenas log
            logger.info(f"Interação recebida (modo simulado): user_id={request.user_id}, content_id={request.content_id}")
        
        return InteractionResponse(
            success=True,
            message="Interação criada com sucesso",
            user_id=request.user_id,
            content_id=request.content_id,
            interaction_type=request.interaction_type
        )
    
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Erro ao criar interação: {e}")
        raise HTTPException(
            status_code=500,
            detail=f"Erro interno ao criar interação: {str(e)}"
        )

def reload_recommendation_model():
    """
    Função auxiliar para recarregar o modelo em background
    """
    try:
        logger.info("Recarregando modelo de recomendação com novas interações...")
        recommendation_model.reload_model()
        logger.info("Modelo de recomendação atualizado com sucesso")
    except Exception as e:
        logger.error(f"Erro ao recarregar modelo: {e}")
