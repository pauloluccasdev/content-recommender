"""
Rotas da API para recomendações
"""
from fastapi import APIRouter, HTTPException
from app.models.recommendation import SimpleRecommendationModel
from app.services.dataset_service import DatasetService
from app.schemas.recommendation import (
    RecommendationRequest,
    RecommendationResponse,
    ContentRecommendation
)
from app.core.config import settings

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
        )
    }
