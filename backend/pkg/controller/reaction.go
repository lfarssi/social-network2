package controller

import (
	"log"
	"real-time-forum/pkg/model"
)

func (s *Server) ReactToPost(request map[string]any) *model.Response {
	response := &model.Response{
		Type:  "reaction",
		Error: "",
	}

	res := s.ValidateSession(request)
	if res.Error != "" {
		response.Error = "Invalid session"
		return response
	}

	postIdRaw, ok := request["postId"]
	if !ok {
		response.Error = "Missing 'postId' field"
		return response
	}
	postId, ok := postIdRaw.(float64)
	if !ok {
		response.Error = "'postId' must be a number"
		return response
	}

	reactionRaw, ok := request["reaction"]
	if !ok {
		response.Error = "Missing 'reaction' field"
		return response
	}

	var reaction *bool = nil

	if reactionRaw != nil {
		reactionVal, ok := reactionRaw.(bool)
		if !ok {
			response.Error = "'reaction' must be a boolean or null"
			return response
		}
		reaction = &reactionVal
	}

	err := s.repository.Reaction().UpsertReaction(res.Userid, int(postId), reaction)
	if err != nil {
		log.Println("Error saving reaction:", err)
		response.Error = "Error saving reaction: " + err.Error()
		return response
	}

	counts, err := s.repository.Reaction().GetReactionCounts(int(postId))
	if err != nil {
		log.Println("Error getting reaction counts:", err)
		response.Error = "Error getting reaction counts: " + err.Error()
		return response
	}

	// Get user's own reaction
	userReaction, err := s.repository.Reaction().GetUserReaction(res.Userid, int(postId))
	if err != nil {
		log.Println("Error getting user reaction:", err)
		response.Error = "Error getting user reaction: " + err.Error()
		return response
	}

	counts.UserReaction = userReaction

	// Create a post object with the reaction counts to return
	post := &model.Post{
		ID:        int(postId),
		Reactions: counts,
	}

	response.Success = true
	response.Posts = []*model.Post{post}

	return response
}
